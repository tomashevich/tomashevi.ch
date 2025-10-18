package middleware

import (
	"fmt"
	"net/http"
	"time"
)

func Cache(duration time.Duration) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			SetCacheRule(w, duration)

			next.ServeHTTP(w, r)
		})
	}
}

func SetCacheRule(w http.ResponseWriter, duration time.Duration) {
	w.Header().Set("Cache-Control", fmt.Sprintf("public, max-age=%d", int64(duration.Seconds())))
}
