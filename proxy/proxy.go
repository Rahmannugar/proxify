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
	upstream := p.loadBalancer.Next()

	upstream.ReverseProxy.ServeHTTP(w, r)
}
