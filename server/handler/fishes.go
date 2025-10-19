package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"
	"tomashevich/server/database"
	"tomashevich/server/middleware"
)

func RegisterFishes(m *http.ServeMux, db *database.Database) {
	listFishes(m, db)
	getFish(m, db)
}

type listFishesResponse struct {
	Seeds []string `json:"seeds"`
}

func listFishes(m *http.ServeMux, db *database.Database) {
	const path = "GET /fishes"
	m.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		pageQuery := r.URL.Query().Get("page")
		page, _ := strconv.ParseInt(pageQuery, 10, 32)
		if page <= 0 {
			page = 0
		}
		page -= 1

		seeds, err := db.GetSeeds(r.Context(), 100, page*100)
		if err != nil {
			http.Error(w, "Failed to get fishes", http.StatusInternalServerError)
			return
		}

		if len(seeds) == 0 {
			seeds = make([]string, 0)
		}

		json.NewEncoder(w).Encode(listFishesResponse{seeds})
	})
}

type getFishResponse struct {
	Seed string `json:"seed"`
}

func getFish(m *http.ServeMux, db *database.Database) {
	const path = "GET /fishes/me"
	m.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		middleware.SetCacheRule(w, time.Hour*168) // week

		id := middleware.GetSoulID(r.Context())
		if id == 0 {
			http.Error(w, "cant get your soul", http.StatusInternalServerError)
			return
		}

		seed, err := db.GetSeed(r.Context(), id)
		if err != nil {
			http.Error(w, "Cant find your soul in fishes", http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(getFishResponse{seed})
	})
}
