package middleware

import (
	"context"
	"net/http"
	"strings"
	"tomashevich/server/database"

	"github.com/google/uuid"
)

type contextKey string

const souldIdKey contextKey = "soulId"

func Helheim(db *database.Database) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip := strings.Split(r.RemoteAddr, ":")[0]
			var id int
			if id, _ = db.GetSoulIDByIP(r.Context(), ip); id == 0 {
				uuid, err := uuid.NewV7()
				if err != nil {
					next.ServeHTTP(w, r)
					return
				}

				db.GiveSoulToHel(r.Context(), uuid.String(), ip)
			}

			if id != 0 {
				r = r.WithContext(context.WithValue(r.Context(), souldIdKey, id))
			}

			next.ServeHTTP(w, r)
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
