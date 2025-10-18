package handler

import (
	"encoding/json"
	"net/http"
	"tomashevich/server/database"
	"tomashevich/server/middleware"
)

func RegisterPixels(m *http.ServeMux, db *database.Database) {
	listPixels(m, db)
	paintPixel(m, db)
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
	m.HandleFunc("GET /pixels", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		pixels, err := db.GetPixels(r.Context())
		if err != nil {
			http.Error(w, "Can get pixels", http.StatusInternalServerError)
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

		json.NewEncoder(w).Encode(listPixelsResponse{allowedColors, colors, x, y})
	})
}

type paintPixelData struct {
	X     int    `json:"x"`
	Y     int    `json:"y"`
	Color string `json:"color"`
}

func paintPixel(m *http.ServeMux, db *database.Database) {
	m.HandleFunc("POST /pixels:paint", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		id := middleware.GetSoulID(r.Context())
		if id == 0 {
			http.Error(w, "cant get your soul", http.StatusInternalServerError)
			return
		}

		var data paintPixelData
		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			http.Error(w, "invalid form", http.StatusUnprocessableEntity)
			return
		}

		defer r.Body.Close()

		color, ok := allowedColors[data.Color]
		if !ok {
			http.Error(w, "invalid color", http.StatusUnprocessableEntity)
			return
		}

		if data.X < 0 || data.Y < 0 {
			http.Error(w, "invalid x/y", http.StatusUnprocessableEntity)
			return
		}

		if err := db.PaintPixel(r.Context(), id, data.X, data.Y, color); err != nil {
			http.Error(w, "cant paint this pixel", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	})
}
