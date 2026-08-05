package config

import (
	"net/http"
	"time"
)

/*
LoadBalancingStrategy represents the algorithm used to distribute
incoming requests across the configured upstreams.
*/
type LoadBalancingStrategy string

const (
	/* RoundRobin distributes requests evenly by cycling through
	all healthy upstreams. */
	RoundRobin LoadBalancingStrategy = "round-robin"

	/* LeastOutstanding distributes requests to the upstream
	with the fewest active requests. */
	LeastOutstanding LoadBalancingStrategy = "least-outstanding"
)

const (
	// DefaultHealthEndpoint is the default endpoint used for health checks.
	DefaultHealthEndpoint = "/health"

	// DefaultHealthInterval is the default interval between health checks.
	DefaultHealthInterval = 10 * time.Second

	// DefaultHealthTimeout is the default timeout for a health check request.
	DefaultHealthTimeout = 2 * time.Second

	// DefaultHealthRetries is the default number of retries before an upstream is considered unhealthy.
	DefaultHealthRetries = 3

	// DefaultExpectedHealthStatus is the expected HTTP status code returned by a healthy upstream.
	DefaultExpectedHealthStatus = http.StatusOK

	// DefaultRequestsPerSecond is the default sustained request rate.
	DefaultRequestsPerSecond = 10

	// DefaultBurst is the default maximum burst size.
	DefaultBurst = 20
)

// Config represents the root application configuration.
type Config struct {
	Server       ServerConfig       `yaml:"server"`
	LoadBalancer LoadBalancerConfig `yaml:"loadBalancer"`
	Health       HealthConfig       `yaml:"health"`
	RateLimit    RateLimitConfig    `yaml:"rateLimit"`
	Upstreams    []Upstream         `yaml:"upstreams"`
}

// ServerConfig defines the proxy server configuration.
type ServerConfig struct {
	Port int `yaml:"port"`
}

/*
LoadBalancerConfig defines how requests are distributed
across upstreams.
*/
type LoadBalancerConfig struct {
	Strategy LoadBalancingStrategy `yaml:"strategy"`
}

// HealthConfig defines the proxy's health check configuration.
type HealthConfig struct {
	Endpoint       string        `yaml:"endpoint"`
	Interval       time.Duration `yaml:"interval"`
	Timeout        time.Duration `yaml:"timeout"`
	Retries        int           `yaml:"retries"`
	ExpectedStatus int           `yaml:"expectedStatus"`
}

// RateLimitConfig defines the proxy's rate limiting configuration.
type RateLimitConfig struct {
	RequestsPerSecond int `yaml:"requestsPerSecond"`
	Burst             int `yaml:"burst"`
}

// Upstream represents a single upstream service.
type Upstream struct {
	ID  string `yaml:"id"`
	URL string `yaml:"url"`
}