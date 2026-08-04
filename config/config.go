package config

// LoadBalancingStrategy represents the algorithm used to distribute
// incoming requests across the configured upstreams.
type LoadBalancingStrategy string

const (
	// RoundRobin distributes requests evenly by cycling through
	// all healthy upstreams.
	RoundRobin LoadBalancingStrategy = "round-robin"

	// LeastOutstanding distributes requests to the upstream
	// with the fewest active requests.
	LeastOutstanding LoadBalancingStrategy = "least-outstanding"
)

// Config represents the root application configuration.
type Config struct {
	Server       ServerConfig       `yaml:"server"`
	LoadBalancer LoadBalancerConfig `yaml:"loadBalancer"`
	Upstreams    []Upstream         `yaml:"upstreams"`
}

// ServerConfig defines the proxy server configuration.
type ServerConfig struct {
	Port int `yaml:"port"`
}

// LoadBalancerConfig defines how requests are distributed
// across upstreams.
type LoadBalancerConfig struct {
	Strategy LoadBalancingStrategy `yaml:"strategy"`
}

// Upstream represents a single upstream service.
type Upstream struct {
	ID  string `yaml:"id"`
	URL string `yaml:"url"`
}