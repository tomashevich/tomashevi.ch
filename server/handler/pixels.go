package handler

import (
	"net/http"
	"time"
	"tomashevich/server/database"
	"tomashevich/server/middleware"
	"tomashevich/utils"
)

func RegisterPixels(m *http.ServeMux, db *database.Database, config *utils.CacheConfig) {
	listPixels(m, db)
	paintPixel(m, db, config)
	registerPixels(m, db)
}

var allowedColors = map[string]int{
	"black":  1,
	"white":  2,
	"red":    3,
	"green":  4,
	"blue":   5,
	"yellow": 6,
	"purple": 7,
	"orange": 8,
}

type listPixelsResponse struct {
	AllowedColors map[string]int `json:"allowed_colors"`
	Colors        []int          `json:"colors"`
	X             []int          `json:"x"`
	Y             []int          `json:"y"`
}

func listPixels(m *http.ServeMux, db *database.Database) {
	const path = "GET /pixels"
	m.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		pixels, err := db.GetPixels(r.Context())
		if err != nil {
			utils.WriteError(w, "Can get pixels", http.StatusInternalServerError)
			return
		}

		if len(pixels) == 0 {
			pixels = make([]database.Pixel, 0)
		}

		capacity := len(pixels)

		colors := make([]int, 0, capacity)
		x := make([]int, 0, capacity)
		y := make([]int, 0, capacity)
		for _, pixel := range pixels {
			colors = append(colors, pixel.Color)
			x = append(x, pixel.X)
			y = append(y, pixel.Y)
		}

		utils.WriteJSON(w, listPixelsResponse{allowedColors, colors, x, y}, http.StatusOK)
	})
}

type paintPixelData struct {
	X     int    `json:"x"`
	Y     int    `json:"y"`
	Color string `json:"color"`
}

func paintPixel(m *http.ServeMux, db *database.Database, config *utils.CacheConfig) {
	const path = "POST /pixels:paint"
	m.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		id := middleware.GetSoulID(r.Context())
		if id == 0 {
			utils.WriteError(w, "cant get your soul", http.StatusInternalServerError)
			return
		}

		soul, err := db.GetSoul(r.Context(), id)
		if err != nil {
			utils.WriteError(w, "cant get your soul", http.StatusInternalServerError)
			return
		}

		var data paintPixelData
		defer r.Body.Close()
		if err := utils.UnmarshalJSON(r.Body, &data); err != nil {
			utils.WriteError(w, "invalid form", http.StatusUnprocessableEntity)
			return
		}

		color, ok := allowedColors[data.Color]
		if !ok {
			utils.WriteError(w, "invalid color", http.StatusUnprocessableEntity)
			return
		}

		if data.X < 0 || data.Y < 0 {
			utils.WriteError(w, "invalid x/y", http.StatusUnprocessableEntity)
			return
		}

		if soul.PaintedPixels >= 10 {
			middleware.SetCacheRule(w, time.Second*time.Duration(config.PixelsLimit)) // dont send again pls
			utils.WriteError(w, "already painted maximum of pixels", http.StatusForbidden)
			return
		}

		if err := db.PaintPixel(r.Context(), soul.Id, data.X, data.Y, color); err != nil {
			utils.WriteError(w, "cant paint this pixel", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	})
}

type registerPixelsData struct {
	Positions []database.PixelPosition `json:"pixels"`
}

func registerPixels(m *http.ServeMux, db *database.Database) {
	const path = "POST /pixels:register"
	m.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		id := middleware.GetSoulID(r.Context())
		if id == 0 {
			utils.WriteError(w, "cant get your soul", http.StatusInternalServerError)
			return
		}

		var data registerPixelsData
		defer r.Body.Close()
		if err := utils.UnmarshalJSON(r.Body, &data); err != nil {
			utils.WriteError(w, "invalid form", http.StatusUnprocessableEntity)
			return
		}

		if ok, _ := db.IsPixelFieldInited(r.Context()); ok {
			utils.WriteError(w, "field already inited", http.StatusAlreadyReported)
			return
		}

		if err := db.InitPixelField(r.Context(), data.Positions, id, allowedColors["white"]); err != nil {
			utils.WriteError(w, "cant init pixel field", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	})

}
