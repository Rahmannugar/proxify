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
	retries        int
	failures       map[string]int
}

// New creates a new health checker.
func New(
	registry *registry.Registry,
	metrics *metrics.Metrics,
	endpoint string,
	interval time.Duration,
	timeout time.Duration,
	expectedStatus int,
	retries int,
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
		retries:        retries,
		failures:       make(map[string]int),
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
		return c.recordFailure(upstream)
	}

	defer resp.Body.Close()

	if resp.StatusCode != c.expectedStatus {
		return c.recordFailure(upstream)
	}

	c.failures[upstream.ID] = 0
	upstream.SetHealthy(true)
	c.logHealth(upstream)

	return true
}

func (c *Checker) recordFailure(upstream *upstream.Upstream) bool {
	c.failures[upstream.ID]++

	if c.failures[upstream.ID] >= c.retries {
		upstream.SetHealthy(false)
	}

	c.logHealth(upstream)

	return upstream.IsHealthy()
}

func (c *Checker) logHealth(upstream *upstream.Upstream) {
	log.Printf(
		"Health check: %s healthy=%t",
		upstream.ID,
		upstream.IsHealthy(),
	)
}
