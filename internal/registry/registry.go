package registry

import (
	"fmt"

	"github.com/rahmannugar/proxify/internal/config"
	"github.com/rahmannugar/proxify/internal/upstream"
)

// Registry manages the collection of configured upstreams.
type Registry struct {
	upstreams []*upstream.Upstream
}

// New creates a new upstream registry.
func New(configUpstreams []config.Upstream) (*Registry, error) {
	upstreams := make([]*upstream.Upstream, 0, len(configUpstreams))

	for _, cfg := range configUpstreams {
		upstream, err := upstream.New(cfg)
		if err != nil {
			return nil, fmt.Errorf(
				"create upstream %q: %w",
				cfg.ID,
				err,
			)
		}

		upstreams = append(upstreams, upstream)
	}

	return &Registry{
		upstreams: upstreams,
	}, nil
}

// List returns all configured upstreams.
func (r *Registry) List() []*upstream.Upstream {
	return r.upstreams
}
