package handler

import (
	"net/http"
	"strconv"
	"time"
	"tomashevich/server/database"
	"tomashevich/server/middleware"
	"tomashevich/server/utils"
)

func RegisterFishes(m *http.ServeMux, db *database.Database, config *utils.CacheConfig) {
	listFishes(m, db)
	getFish(m, db, config)
}

type listFishesResponse struct {
	Seeds []string `json:"seeds"`
}

func listFishes(m *http.ServeMux, db *database.Database) {
	const path = "GET /fishes"
	m.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		pageQuery := r.URL.Query().Get("page")
		page, _ := strconv.ParseInt(pageQuery, 10, 32)
		if page <= 0 {
			page = 0
		}
		page -= 1

		seeds, err := db.GetSeeds(r.Context(), 100, page*100)
		if err != nil {
			utils.WriteError(w, "Failed to get fishes", http.StatusInternalServerError)
			return
		}

		if len(seeds) == 0 {
			seeds = make([]string, 0)
		}

		utils.WriteJSON(w, listFishesResponse{seeds}, http.StatusOK)
	})
}

type getFishResponse struct {
	Seed string `json:"seed"`
}

func getFish(m *http.ServeMux, db *database.Database, config *utils.CacheConfig) {
	const path = "GET /fishes/me"
	m.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		middleware.SetCacheRule(w, time.Second*time.Duration(config.FishesMe)) // week

		id := middleware.GetSoulID(r.Context())
		if id == 0 {
			utils.WriteError(w, "cant get your soul", http.StatusInternalServerError)
			return
		}

		seed, err := db.GetSeed(r.Context(), id)
		if err != nil {
			utils.WriteError(w, "cant find your soul in fishes", http.StatusInternalServerError)
			return
		}

		utils.WriteJSON(w, getFishResponse{seed}, http.StatusOK)
	})
}
