package upstream

import (
	"net/http/httputil"
	"net/url"

	"github.com/rahmannugar/proxify/config"
)

// Upstream represents a runtime upstream service.
type Upstream struct {
	ID           string
	URL          string
	ReverseProxy *httputil.ReverseProxy
}

// New creates a runtime upstream from configuration.
func New(cfg config.Upstream) (*Upstream, error) {
	target, err := url.Parse(cfg.URL)
	if err != nil {
		return nil, err
	}

	return &Upstream{
		ID:           cfg.ID,
		URL:          cfg.URL,
		ReverseProxy: httputil.NewSingleHostReverseProxy(target),
	}, nil
}
