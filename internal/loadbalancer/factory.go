package loadbalancer

import (
	"fmt"

	"github.com/rahmannugar/proxify/internal/config"
	"github.com/rahmannugar/proxify/internal/registry"
)

/*
Creates a load balancer for the configured strategy.
*/
func NewLoadBalancer(
	strategy config.LoadBalancingStrategy,
	registry *registry.Registry,
) (LoadBalancer, error) {
	switch strategy {
	case config.RoundRobin:
		return NewRoundRobin(registry), nil

	case config.LeastOutstanding:
		return NewLeastOutstanding(registry), nil

	default:
		return nil, fmt.Errorf(
			"unsupported load balancing strategy: %q",
			strategy,
		)
	}
}
