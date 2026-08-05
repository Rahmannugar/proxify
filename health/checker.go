package health

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/rahmannugar/proxify/registry"
)

// Checker periodically verifies the health of upstream services.
type Checker struct {
	client         *http.Client
	registry       *registry.Registry
	interval       time.Duration
	endpoint       string
	expectedStatus int
}

// New creates a new health checker.
func New(
	registry *registry.Registry,
	endpoint string,
	interval time.Duration,
	timeout time.Duration,
	expectedStatus int,
) *Checker {
	return &Checker{
		client: &http.Client{
			Timeout: timeout,
		},
		registry:       registry,
		interval:       interval,
		endpoint:       endpoint,
		expectedStatus: expectedStatus,
	}
}

// Start begins periodically checking upstream health.
func (c *Checker) Start(ctx context.Context) {
	ticker := time.NewTicker(c.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return

		case <-ticker.C:
			c.check()
		}
	}
}

// check performs a single health check against every upstream.
func (c *Checker) check() {
	for _, upstream := range c.registry.List() {
		resp, err := c.client.Get(upstream.URL + c.endpoint)
		if err != nil {
			upstream.SetHealthy(false)

			log.Printf(
				"Health check: %s healthy=%t",
				upstream.ID,
				upstream.IsHealthy(),
			)

			continue
		}

		resp.Body.Close()

		upstream.SetHealthy(resp.StatusCode == c.expectedStatus)

		log.Printf(
			"Health check: %s healthy=%t",
			upstream.ID,
			upstream.IsHealthy(),
		)
	}
}
