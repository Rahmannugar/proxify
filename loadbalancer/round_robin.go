package loadbalancer

import (
	"sync"

	"github.com/rahmannugar/proxify/registry"
	"github.com/rahmannugar/proxify/upstream"
)

// RoundRobin distributes requests across upstreams.
type RoundRobin struct {
	registry *registry.Registry
	mu       sync.Mutex
	index    int
}

// NewRoundRobin creates a new round-robin load balancer.
func NewRoundRobin(registry *registry.Registry) *RoundRobin {
	return &RoundRobin{
		registry: registry,
	}
}

// Next returns the next healthy upstream in the rotation.
func (r *RoundRobin) Next() (*upstream.Upstream, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	upstreams := r.registry.List()

	for i := 0; i < len(upstreams); i++ {
		upstream := upstreams[r.index]

		r.index = (r.index + 1) % len(upstreams)

		if upstream.IsHealthy() {
			return upstream, nil
		}
	}

	return nil, ErrNoHealthyUpstreams
}
