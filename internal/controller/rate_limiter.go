package controller

import (
	"net/http"

	"golang.org/x/time/rate"
)

const limit = 50
const burst = 100

func rateLimitMiddleware(limiters userLimiters, next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, ok := limiters[r.RemoteAddr]
		if !ok {
			limiters[r.RemoteAddr] = rate.NewLimiter(limit, burst)
		}

		if !limiters[r.RemoteAddr].Allow() {
			http.Error(w, http.StatusText(http.StatusTooManyRequests), http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}
