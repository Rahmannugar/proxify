package config

import (
	"fmt"
	"net/url"
)

// Validate ensures the application configuration is valid before the proxy starts.
func Validate(cfg *Config) error {
	if cfg.Server.Port <= 0 || cfg.Server.Port > 65535 {
		return fmt.Errorf("invalid server port: %d", cfg.Server.Port)
	}

	if len(cfg.Upstreams) == 0 {
		return fmt.Errorf("at least one upstream must be configured")
	}

	switch cfg.LoadBalancer.Strategy {
	case RoundRobin, LeastOutstanding:
	default:
		return fmt.Errorf(
			"unsupported load balancing strategy: %q",
			cfg.LoadBalancer.Strategy,
		)
	}

	if cfg.Health.Endpoint == "" {
		cfg.Health.Endpoint = DefaultHealthEndpoint
	}

	if cfg.Health.Interval <= 0 {
		cfg.Health.Interval = DefaultHealthInterval
	}

	if cfg.Health.Timeout <= 0 {
		cfg.Health.Timeout = DefaultHealthTimeout
	}

	if cfg.Health.Retries <= 0 {
		cfg.Health.Retries = DefaultHealthRetries
	}

	if cfg.Health.ExpectedStatus == 0 {
		cfg.Health.ExpectedStatus = DefaultExpectedHealthStatus
	}

	seen := make(map[string]struct{})

	for _, upstream := range cfg.Upstreams {
		if upstream.ID == "" {
			return fmt.Errorf("upstream id cannot be empty")
		}

		if _, exists := seen[upstream.ID]; exists {
			return fmt.Errorf("duplicate upstream id: %q", upstream.ID)
		}

		seen[upstream.ID] = struct{}{}

		if _, err := url.ParseRequestURI(upstream.URL); err != nil {
			return fmt.Errorf(
				"invalid upstream URL %q: %w",
				upstream.URL,
				err,
			)
		}
	}

	return nil
}
