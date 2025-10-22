package middleware

import (
	"context"
	"net/http"
	"tomashevich/server/database"
	"tomashevich/server/utils"

	"github.com/google/uuid"
)

type contextKey string

const souldIdKey contextKey = "soulId"

func Helheim(db *database.Database) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip := utils.GetIPAddr(r)
			var id int
			if id, _ = db.GetSoulIDByIP(r.Context(), ip); id == 0 {
				uuid, err := uuid.NewV7()
				if err != nil {
					next.ServeHTTP(w, r)
					return
				}

				if id, err = db.GiveSoulToHel(r.Context(), uuid.String(), ip); err != nil {
					next.ServeHTTP(w, r)
					return
				}
			}

			next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), souldIdKey, id)))
		})
	}
}

func GetSoulID(ctx context.Context) int {
	id, ok := ctx.Value(souldIdKey).(int)
	if !ok {
		return 0
	}
	return id
}
