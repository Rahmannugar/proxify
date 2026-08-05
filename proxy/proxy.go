package proxy

import (
	"net/http"

	"github.com/rahmannugar/proxify/loadbalancer"
)

// Proxy forwards incoming HTTP requests to upstream services.
type Proxy struct {
	loadBalancer loadbalancer.LoadBalancer
}

// New creates a new reverse proxy.
func New(lb loadbalancer.LoadBalancer) *Proxy {
	return &Proxy{
		loadBalancer: lb,
	}
}

// ServeHTTP forwards the incoming request to the selected upstream.
func (p *Proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	upstream, err := p.loadBalancer.Next()
	if err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusServiceUnavailable,
		)

		return
	}

	upstream.ReverseProxy.ServeHTTP(w, r)
}
