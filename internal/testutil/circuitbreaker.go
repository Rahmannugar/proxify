package testutil

import (
	"time"

	"github.com/rahmannugar/proxify/internal/config"
)

// CircuitBreakerConfig returns a reusable circuit breaker configuration for tests.
func CircuitBreakerConfig() config.CircuitBreakerConfig {
	return config.CircuitBreakerConfig{
		FailureThreshold: 5,
		RecoveryTimeout:  30 * time.Second,
		SuccessThreshold: 2,
	}
}
