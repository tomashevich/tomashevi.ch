package utils

import (
	"encoding/json"
	"net/http"
)

type ResponseError struct {
	Details string `json:"details"`
}

func WriteError(w http.ResponseWriter, details string, statusCode int) {
	w.WriteHeader(statusCode)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ResponseError{details})
}

func WriteJSON(w http.ResponseWriter, v any, statusCode int) {
	w.WriteHeader(statusCode)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(v); err != nil {
		WriteError(w, "cant encode result", http.StatusInternalServerError)
	}
}
