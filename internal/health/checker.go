package health

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/rahmannugar/proxify/internal/metrics"
	"github.com/rahmannugar/proxify/internal/registry"
	"github.com/rahmannugar/proxify/internal/upstream"
)

// Checker periodically verifies the health of upstream services.
type Checker struct {
	client         *http.Client
	registry       *registry.Registry
	metrics        *metrics.Metrics
	interval       time.Duration
	endpoint       string
	expectedStatus int
}

// New creates a new health checker.
func New(
	registry *registry.Registry,
	metrics *metrics.Metrics,
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
		metrics:        metrics,
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

// check performs a health check against every configured upstream.
func (c *Checker) check() {
	healthy := 0
	unhealthy := 0

	for _, upstream := range c.registry.List() {
		if c.checkUpstream(upstream) {
			healthy++
		} else {
			unhealthy++
		}
	}

	c.metrics.SetHealthyUpstreams(healthy)
	c.metrics.SetUnhealthyUpstreams(unhealthy)
}

// checkUpstream performs a health check against a single upstream.
func (c *Checker) checkUpstream(upstream *upstream.Upstream) bool {
	resp, err := c.client.Get(upstream.URL + c.endpoint)
	if err != nil {
		upstream.SetHealthy(false)

		log.Printf(
			"Health check: %s healthy=%t",
			upstream.ID,
			upstream.IsHealthy(),
		)

		return false
	}

	defer resp.Body.Close()

	healthy := resp.StatusCode == c.expectedStatus

	upstream.SetHealthy(healthy)

	log.Printf(
		"Health check: %s healthy=%t",
		upstream.ID,
		upstream.IsHealthy(),
	)

	return healthy
}
