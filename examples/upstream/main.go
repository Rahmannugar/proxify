package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"
)

func main() {
	name := flag.String(
		"name",
		"api-1",
		"Upstream instance name.",
	)

	port := flag.Int(
		"port",
		9001,
		"Port to listen on.",
	)

	delay := flag.Duration(
		"delay",
		0,
		"Artificial response delay.",
	)

	flag.Parse()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(*delay)

		if err := json.NewEncoder(w).Encode(map[string]string{
			"instance": *name,
		}); err != nil {
			http.Error(
				w,
				"failed to encode response",
				http.StatusInternalServerError,
			)

			return
		}
	})

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	addr := fmt.Sprintf(":%d", *port)

	log.Printf(
		"Starting %s on %s (delay=%s)",
		*name,
		addr,
		delay.String(),
	)

	log.Fatal(http.ListenAndServe(addr, nil))
}
