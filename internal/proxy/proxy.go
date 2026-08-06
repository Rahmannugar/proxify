package proxy

import (
	"net/http"

	"github.com/rahmannugar/proxify/internal/loadbalancer"
	"github.com/rahmannugar/proxify/internal/metrics"
)

// upstreamRecorder records the upstream that handled a request.
type upstreamRecorder interface {
	SetUpstream(string)
}

// Proxy forwards incoming HTTP requests to upstream services.
type Proxy struct {
	loadBalancer loadbalancer.LoadBalancer
	metrics      *metrics.Metrics
}

// New creates a new reverse proxy.
func New(
	lb loadbalancer.LoadBalancer,
	metrics *metrics.Metrics,
) *Proxy {
	return &Proxy{
		loadBalancer: lb,
		metrics:      metrics,
	}
}

// ServeHTTP forwards the incoming request to the selected upstream.
func (p *Proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	p.metrics.IncrementRequests()
	p.metrics.IncrementInFlight()
	defer p.metrics.DecrementInFlight()

	upstream, err := p.loadBalancer.Next()
	if err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusServiceUnavailable,
		)
		return
	}

	if !upstream.Breaker.Allow() {
		http.Error(
			w,
			"upstream circuit breaker is open",
			http.StatusServiceUnavailable,
		)
		return
	}

	upstream.IncrementActiveRequests()
	defer upstream.DecrementActiveRequests()

	if recorder, ok := w.(upstreamRecorder); ok {
		recorder.SetUpstream(upstream.ID)
	}

	upstream.ReverseProxy.ServeHTTP(w, r)
}
