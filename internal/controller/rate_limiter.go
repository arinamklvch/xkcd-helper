package controller

import (
	"net"
	"net/http"

	"golang.org/x/time/rate"
)

func rateLimit(limiters *userLimiters, limit rate.Limit, burst int, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		host, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			host = r.RemoteAddr
		}

		limiters.mu.Lock()
		limiter, ok := limiters.limiters[host]
		if !ok {
			limiter = rate.NewLimiter(limit, burst)
			limiters.limiters[host] = limiter
		}
		limiters.mu.Unlock()

		if !limiter.Allow() {
			http.Error(w, http.StatusText(http.StatusTooManyRequests), http.StatusTooManyRequests)
			return
		}

		next(w, r)
	}
}
