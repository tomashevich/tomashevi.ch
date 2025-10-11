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
		if pageQuery == "" {
			pageQuery = "1"
		}

		page, err := strconv.ParseInt(pageQuery, 10, 32)
		if err != nil {
			http.Error(w, "Page param must be int", http.StatusUnprocessableEntity)
			return
		}
		if page <= 0 {
			http.Error(w, "Page param must be positive", http.StatusUnprocessableEntity)
			return
		}
		page -= 1

		fishes, err := db.GetFishes(r.Context(), 100, page*100)
		if err != nil {
			http.Error(w, "Failed to get fishes", http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(fishes)
	})
}

func getFish(m *http.ServeMux, db *database.Database) {
	m.HandleFunc("GET /fishes/me", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		fish, err := db.GetFishByIP(r.Context(), strings.Split(r.RemoteAddr, ":")[0])
		if err != nil {
			http.Error(w, "Cant find your soul in fishes", http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(fish)
	})
}
