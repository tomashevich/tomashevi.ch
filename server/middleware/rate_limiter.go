package middleware

import (
	"net/http"
	"strconv"
	"sync"
	"time"
	"tomashevich/server/utils"
)

type RequestLimitState struct {
	Count     int
	ResetTime time.Time
}

type RateLimiter struct {
	mu          sync.Mutex
	requests    map[string]*RequestLimitState
	MaxRequests int
	InDuration  time.Duration
}

func NewRateLimiter(maxRequests int, inDuration time.Duration) *RateLimiter {
	return &RateLimiter{
		requests:    make(map[string]*RequestLimitState),
		MaxRequests: maxRequests,
		InDuration:  inDuration,
	}
}

func (rl *RateLimiter) RunCleaner() {
	ticker := time.NewTicker(rl.InDuration)

	go func() {
		for range ticker.C {
			rl.mu.Lock()

			now := time.Now()
			for ip, state := range rl.requests {
				if now.After(state.ResetTime) {
					delete(rl.requests, ip)
				}
			}

			rl.mu.Unlock()
		}
	}()
}

func (rl *RateLimiter) Middleware(isBehindProxy bool) func(next http.Handler) http.Handler {
	rl.RunCleaner()
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip := utils.GetIPAddr(r, isBehindProxy)
			now := time.Now()

			rl.mu.Lock()
			defer rl.mu.Unlock()

			state, exists := rl.requests[ip]
			if exists && now.After(state.ResetTime) {
				delete(rl.requests, ip)
				exists = false // Обрабатываем как новый запрос
			}

			if !exists {
				state = &RequestLimitState{
					1,
					now.Add(rl.InDuration),
				}
				rl.requests[ip] = state
				rl.SetHeaders(w, state)
				next.ServeHTTP(w, r)
				return
			}

			rl.SetHeaders(w, state)

			if state.Count >= rl.MaxRequests {
				utils.WriteError(w, "rate limit", http.StatusTooManyRequests)
				return
			}

			state.Count++
			next.ServeHTTP(w, r)
		})
	}
}

func (rl *RateLimiter) SetHeaders(w http.ResponseWriter, state *RequestLimitState) {
	w.Header().Set("X-RateLimit-Limit", strconv.Itoa(rl.MaxRequests))
	w.Header().Set("X-RateLimit-Remaining", strconv.Itoa(rl.MaxRequests-state.Count))
	w.Header().Set("X-RateLimit-Reset", state.ResetTime.Format(time.RFC1123))
}
