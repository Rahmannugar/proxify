package loadbalancer

import (
	"github.com/rahmannugar/proxify/config"
	"github.com/rahmannugar/proxify/registry"
)

// RoundRobin distributes requests across upstreams.
type RoundRobin struct {
	registry *registry.Registry
	index    int
}

// NewRoundRobin creates a new round-robin load balancer.
func NewRoundRobin(registry *registry.Registry) *RoundRobin {
	return &RoundRobin{
		registry: registry,
	}
}

// Next returns the next upstream.
func (r *RoundRobin) Next() config.Upstream {
	upstreams := r.registry.List()

	upstream := upstreams[r.index]

	r.index = (r.index + 1) % len(upstreams)

	return upstream
}
