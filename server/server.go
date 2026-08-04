package server

import (
	"fmt"
	"net/http"

	"github.com/rahmannugar/proxify/config"
)

// Server represents the HTTP server that accepts incoming
// client requests.
type Server struct {
	httpServer *http.Server
}

// New creates a new HTTP server.
func New(cfg *config.Config, handler http.Handler) *Server {
	return &Server{
		httpServer: &http.Server{
			Addr:    fmt.Sprintf(":%d", cfg.Server.Port),
			Handler: handler,
		},
	}
}

// Start begins listening for incoming HTTP requests.
func (s *Server) Start() error {
	return s.httpServer.ListenAndServe()
}
