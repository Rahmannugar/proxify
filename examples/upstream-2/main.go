package main

import (
	"encoding/json"
	"log"
	"net/http"
)

// Response represents the response returned by the example upstream.
type Response struct {
	Instance string `json:"instance"`
}

func main() {
	mux := http.NewServeMux()

	// Health endpoint used by Proxify to determine whether
	// the upstream is healthy.
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Root endpoint returns the identity of the upstream.
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		response := Response{
			Instance: "api-2",
		}

		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(
				w,
				"failed to encode response",
				http.StatusInternalServerError,
			)

			return
		}
	})

	log.Println("Example upstream listening on :9002")

	log.Fatal(http.ListenAndServe(":9002", mux))
}
