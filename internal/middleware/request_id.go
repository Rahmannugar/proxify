package middleware

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"

	"github.com/rahmannugar/proxify/internal/requestid"
)

const RequestIDHeader = "X-Request-ID"

// RequestID ensures every request has a unique request ID.
func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := r.Header.Get(RequestIDHeader)

		if requestID == "" {
			var err error

			requestID, err = generateRequestID()
			if err != nil {
				http.Error(
					w,
					"failed to generate request ID",
					http.StatusInternalServerError,
				)

				return
			}
		}

		// Return the request ID to the client.
		w.Header().Set(RequestIDHeader, requestID)

		// Store the request ID in the request context.
		ctx := requestid.NewContext(
			r.Context(),
			requestID,
		)

		next.ServeHTTP(
			w,
			r.WithContext(ctx),
		)
	})
}

// generateRequestID creates a cryptographically random request ID.
func generateRequestID() (string, error) {
	var bytes [16]byte

	if _, err := rand.Read(bytes[:]); err != nil {
		return "", err
	}

	return hex.EncodeToString(bytes[:]), nil
}
