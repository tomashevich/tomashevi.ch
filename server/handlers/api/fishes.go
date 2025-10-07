package api

import (
	"net/http"
)

func RegisterFishes(m *http.ServeMux) {
	listFishes(m)
	getFish(m)
}

func listFishes(m *http.ServeMux) {
	m.HandleFunc("GET /fishes", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("fishes 1 2 3"))
	})
}

func getFish(m *http.ServeMux) {
	m.HandleFunc("GET /fishes/{id}", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("fish with id " + r.PathValue("id")))
	})
}
