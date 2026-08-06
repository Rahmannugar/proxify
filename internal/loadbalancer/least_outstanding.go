package loadbalancer

import (
	"fmt"

	"github.com/rahmannugar/proxify/internal/registry"
	"github.com/rahmannugar/proxify/internal/upstream"
)

// LeastOutstanding selects the upstream with the fewest active requests.
type LeastOutstanding struct {
	registry *registry.Registry
}

// NewLeastOutstanding creates a least outstanding load balancer.
func NewLeastOutstanding(
	registry *registry.Registry,
) *LeastOutstanding {
	return &LeastOutstanding{
		registry: registry,
	}
}

// Next returns the least busy healthy upstream.
func (l *LeastOutstanding) Next() (*upstream.Upstream, error) {
	upstreams := l.registry.List()

	var selected *upstream.Upstream

	for _, upstream := range upstreams {
		if !upstream.IsHealthy() {
			continue
		}

		if selected == nil ||
			upstream.ActiveRequests() < selected.ActiveRequests() {
			selected = upstream
		}
	}

	if selected == nil {
		return nil, fmt.Errorf("no healthy upstreams available")
	}

	return selected, nil
}
