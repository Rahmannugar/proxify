package server

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rahmannugar/proxify/internal/config"
)

const shutdownTimeout = 10 * time.Second

// Server represents the HTTP server that accepts incoming client requests.
type Server struct {
	httpServer *http.Server
}

// New creates a new HTTP server.
func New(
	cfg *config.Config,
	proxyHandler http.Handler,
	metricsHandler http.Handler,
) *Server {
	mux := http.NewServeMux()

	mux.Handle("/", proxyHandler)
	mux.Handle("/metrics", metricsHandler)

	return &Server{
		httpServer: &http.Server{
			Addr:    fmt.Sprintf(":%d", cfg.Server.Port),
			Handler: mux,
		},
	}
}

// Start begins listening for incoming HTTP requests.
func (s *Server) Start() error {
	/*
		ListenAndServe blocks, so run it in a separate goroutine
		while we wait for either a server error or an OS shutdown signal.
	*/
	serverErr := make(chan error, 1)

	go func() {
		if err := s.httpServer.ListenAndServe(); err != nil &&
			!errors.Is(err, http.ErrServerClosed) {
			serverErr <- err
		}
	}()

	/*
		Wait for either:
		- an unexpected server error
		- a shutdown signal from the operating system
	*/
	stop := make(chan os.Signal, 1)

	signal.Notify(
		stop,
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer signal.Stop(stop)

	select {
	case err := <-serverErr:
		return err

	case <-stop:
		log.Println("Received shutdown signal. Shutting down server...")

		ctx, cancel := context.WithTimeout(
			context.Background(),
			shutdownTimeout,
		)
		defer cancel()

		if err := s.httpServer.Shutdown(ctx); err != nil {
			return err
		}

		log.Println("Server shut down gracefully.")

		return nil
	}
}
