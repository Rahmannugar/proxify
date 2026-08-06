package metrics

import "sync/atomic"

// Snapshot represents the current operational metrics.
type Snapshot struct {
	RequestsTotal      int64
	RequestsInFlight   int64
	RateLimited        int64
	HealthyUpstreams   int64
	UnhealthyUpstreams int64
}

// Metrics stores operational metrics for the proxy.
type Metrics struct {
	requestsTotal      atomic.Int64
	requestsInFlight   atomic.Int64
	rateLimited        atomic.Int64
	healthyUpstreams   atomic.Int64
	unhealthyUpstreams atomic.Int64
}

// New creates a new metrics collector.
func New() *Metrics {
	return &Metrics{}
}

// IncrementRequests increments the total requests counter.
func (m *Metrics) IncrementRequests() {
	m.requestsTotal.Add(1)
}

// IncrementInFlight increments the in-flight requests counter.
func (m *Metrics) IncrementInFlight() {
	m.requestsInFlight.Add(1)
}

// DecrementInFlight decrements the in-flight requests counter.
func (m *Metrics) DecrementInFlight() {
	m.requestsInFlight.Add(-1)
}

// IncrementRateLimited increments the rate-limited requests counter.
func (m *Metrics) IncrementRateLimited() {
	m.rateLimited.Add(1)
}

// SetHealthyUpstreams updates the number of healthy upstreams.
func (m *Metrics) SetHealthyUpstreams(count int) {
	m.healthyUpstreams.Store(int64(count))
}

// SetUnhealthyUpstreams updates the number of unhealthy upstreams.
func (m *Metrics) SetUnhealthyUpstreams(count int) {
	m.unhealthyUpstreams.Store(int64(count))
}

// Snapshot returns the current operational metrics.
func (m *Metrics) Snapshot() Snapshot {
	return Snapshot{
		RequestsTotal:      m.requestsTotal.Load(),
		RequestsInFlight:   m.requestsInFlight.Load(),
		RateLimited:        m.rateLimited.Load(),
		HealthyUpstreams:   m.healthyUpstreams.Load(),
		UnhealthyUpstreams: m.unhealthyUpstreams.Load(),
	}
}
