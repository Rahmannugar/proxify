package transport

import (
	"net"
	"net/http"

	"github.com/rahmannugar/proxify/internal/config"
)

// New creates a shared HTTP transport for communicating with upstream services.
func New(cfg config.TransportConfig) *http.Transport {
	return &http.Transport{
		DialContext: (&net.Dialer{
			Timeout: cfg.Dial,
		}).DialContext,

		ResponseHeaderTimeout: cfg.ResponseHeader,
		IdleConnTimeout:       cfg.IdleConn,
		TLSHandshakeTimeout:   cfg.TLSHandshake,

		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 10,
	}
}
