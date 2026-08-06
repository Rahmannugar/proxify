package middleware

import (
	"log"
	"net/http"
	"time"
)

/*
Logging logs information about each HTTP request after
the downstream handler has finished processing it.
*/
func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		rw := &responseWriter{
			ResponseWriter: w,
		}

		next.ServeHTTP(rw, r)

		if rw.upstreamID == "" {
			rw.upstreamID = "-"
		}

		log.Printf(
			"%s %s %d %s upstream=%s",
			r.Method,
			r.URL.Path,
			rw.statusCode,
			time.Since(start),
			rw.upstreamID,
		)
	})
}
