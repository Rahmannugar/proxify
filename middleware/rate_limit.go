package middleware

import (
	"log"
	"net"
	"net/http"

	"github.com/rahmannugar/proxify/ratelimiter"
)

// RateLimit limits incoming requests on a per-client basis.
func RateLimit(
	limiter *ratelimiter.Limiter,
	next http.Handler,
) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		clientIP, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			http.Error(
				w,
				"Failed to determine client address.",
				http.StatusInternalServerError,
			)

			return
		}

		if !limiter.Allow(clientIP) {
			log.Printf(
				"Rate limit exceeded: client=%s",
				clientIP,
			)

			http.Error(
				w,
				http.StatusText(http.StatusTooManyRequests),
				http.StatusTooManyRequests,
			)

			return
		}

		next.ServeHTTP(w, r)
	})
}