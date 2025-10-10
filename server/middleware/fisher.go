package middleware

import (
	"net/http"
	"strings"
	"tomashevich/server/database"

	"github.com/google/uuid"
)

func Fisher(db *database.Database) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip := strings.Split(r.RemoteAddr, ":")[0]
			if fish, _ := db.GetFishByIP(r.Context(), ip); fish.Address == "" {
				uuid, err := uuid.NewV7()
				if err != nil {
					next.ServeHTTP(w, r)
					return
				}

				db.AddFish(r.Context(), uuid.String(), ip)
			}

			next.ServeHTTP(w, r)
		})
	}
}
