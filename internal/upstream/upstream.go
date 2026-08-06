package upstream

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
	"sync/atomic"

	"github.com/rahmannugar/proxify/internal/circuitbreaker"
	"github.com/rahmannugar/proxify/internal/config"
)

// Upstream represents a runtime upstream service.
type Upstream struct {
	mu sync.RWMutex

	ID             string
	URL            string
	ReverseProxy   *httputil.ReverseProxy
	Breaker        *circuitbreaker.Breaker
	healthy        bool
	activeRequests atomic.Int64
}

// IsHealthy reports whether the upstream is healthy.
func (u *Upstream) IsHealthy() bool {
	u.mu.RLock()
	defer u.mu.RUnlock()

	return u.healthy
}

// SetHealthy updates the upstream health status.
func (u *Upstream) SetHealthy(healthy bool) {
	u.mu.Lock()
	defer u.mu.Unlock()

	u.healthy = healthy
}

// IncrementActiveRequests increases the number of active requests.
func (u *Upstream) IncrementActiveRequests() {
	u.activeRequests.Add(1)
}

// DecrementActiveRequests decreases the number of active requests.
func (u *Upstream) DecrementActiveRequests() {
	u.activeRequests.Add(-1)
}

// ActiveRequests returns the number of active requests.
func (u *Upstream) ActiveRequests() int64 {
	return u.activeRequests.Load()
}

// New creates a runtime upstream from configuration.
func New(
	cfg config.Upstream,
	cbCfg config.CircuitBreakerConfig,
	transport http.RoundTripper,
) (*Upstream, error) {
	target, err := url.Parse(cfg.URL)
	if err != nil {
		return nil, err
	}

	proxy := httputil.NewSingleHostReverseProxy(target)
	proxy.Transport = transport

	return &Upstream{
		ID:           cfg.ID,
		URL:          cfg.URL,
		ReverseProxy: proxy,
		Breaker: circuitbreaker.New(
			cbCfg.FailureThreshold,
			cbCfg.SuccessThreshold,
			cbCfg.RecoveryTimeout,
		),
		healthy: true,
	}, nil
}
