package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"tomashevich/server/database"
)

func RegisterFishes(m *http.ServeMux, db *database.Database) {
	listFishes(m, db)
	getFish(m, db)
}

func listFishes(m *http.ServeMux, db *database.Database) {
	m.HandleFunc("GET /fishes", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		pageQuery := r.URL.Query().Get("page")
		page, _ := strconv.ParseInt(pageQuery, 10, 32)
		if page <= 0 {
			http.Error(w, "Page param must be uint", http.StatusUnprocessableEntity)
			return
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

		json.NewEncoder(w).Encode(seeds)
	})
}

func getFish(m *http.ServeMux, db *database.Database) {
	m.HandleFunc("GET /fishes/me", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		seed, err := db.GetSeedByIP(r.Context(), strings.Split(r.RemoteAddr, ":")[0])
		if err != nil {
			http.Error(w, "Cant find your soul in fishes", http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(seed)
	})
}
