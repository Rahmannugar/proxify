package loadbalancer

import "github.com/rahmannugar/proxify/internal/upstream"

// LoadBalancer selects an upstream to receive the next request.
type LoadBalancer interface {
	Next() (*upstream.Upstream, error)
}
