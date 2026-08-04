package proxy

import (
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/rahmannugar/proxify/config"
)

// Proxy forwards incoming HTTP requests to an upstream service.
type Proxy struct {
	reverseProxy *httputil.ReverseProxy
}

// New creates a new reverse proxy.
//
// For now, all requests are forwarded to the first configured upstream.
func New(cfg *config.Config) (*Proxy, error) {
	target, err := url.Parse(cfg.Upstreams[0].URL)
	if err != nil {
		return nil, err
	}

	return &Proxy{
		reverseProxy: httputil.NewSingleHostReverseProxy(target),
	}, nil
}

// ServeHTTP forwards the incoming request to the configured upstream.
func (p *Proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	p.reverseProxy.ServeHTTP(w, r)
}
