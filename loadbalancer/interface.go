package loadbalancer

import "github.com/rahmannugar/proxify/config"

// LoadBalancer selects an upstream to receive the next request.
type LoadBalancer interface {
	Next() config.Upstream
}
