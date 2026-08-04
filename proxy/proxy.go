package proxy

import (
	"net/http"
	"net/http/httputil"
	"net/url"

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

	target, err := url.Parse(upstream.URL)
	if err != nil {
		http.Error(
			w,
			"failed to resolve upstream",
			http.StatusInternalServerError,
		)

		return
	}

	reverseProxy := httputil.NewSingleHostReverseProxy(target)

	reverseProxy.ServeHTTP(w, r)
}
