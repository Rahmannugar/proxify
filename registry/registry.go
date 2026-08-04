package registry

import "github.com/rahmannugar/proxify/config"

// Registry manages the collection of configured upstreams.
type Registry struct {
	upstreams []config.Upstream
}

// New creates a new upstream registry.
func New(upstreams []config.Upstream) *Registry {
	return &Registry{
		upstreams: upstreams,
	}
}

// List returns all configured upstreams.
func (r *Registry) List() []config.Upstream {
	return r.upstreams
}
