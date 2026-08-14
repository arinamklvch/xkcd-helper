package controller

import (
	"net/http"

	"golang.org/x/time/rate"
)

func rateLimit(limiters userLimiters, limit rate.Limit, burst int, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, ok := limiters[r.RemoteAddr]
		if !ok {
			limiters[r.RemoteAddr] = rate.NewLimiter(limit, burst)
		}

		if !limiters[r.RemoteAddr].Allow() {
			http.Error(w, http.StatusText(http.StatusTooManyRequests), http.StatusTooManyRequests)
			return
		}

		next(w, r)
	}
}
