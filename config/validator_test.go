package config

import "testing"

func TestValidate(t *testing.T) {
	tests := []struct {
		name    string
		config  Config
		wantErr bool
	}{
		{
			name: "valid configuration",
			config: Config{
				Server: ServerConfig{
					Port: 8080,
				},
				LoadBalancer: LoadBalancerConfig{
					Strategy: RoundRobin,
				},
				Upstreams: []Upstream{
					{
						ID:  "api-1",
						URL: "http://localhost:9001",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "invalid server port",
			config: Config{
				Server: ServerConfig{
					Port: 0,
				},
				LoadBalancer: LoadBalancerConfig{
					Strategy: RoundRobin,
				},
				Upstreams: []Upstream{
					{
						ID:  "api-1",
						URL: "http://localhost:9001",
					},
				},
			},
			wantErr: true,
		},
		{
			name: "no upstreams configured",
			config: Config{
				Server: ServerConfig{
					Port: 8080,
				},
				LoadBalancer: LoadBalancerConfig{
					Strategy: RoundRobin,
				},
			},
			wantErr: true,
		},
		{
			name: "unsupported load balancing strategy",
			config: Config{
				Server: ServerConfig{
					Port: 8080,
				},
				LoadBalancer: LoadBalancerConfig{
					Strategy: "random",
				},
				Upstreams: []Upstream{
					{
						ID:  "api-1",
						URL: "http://localhost:9001",
					},
				},
			},
			wantErr: true,
		},
		{
			name: "duplicate upstream ids",
			config: Config{
				Server: ServerConfig{
					Port: 8080,
				},
				LoadBalancer: LoadBalancerConfig{
					Strategy: RoundRobin,
				},
				Upstreams: []Upstream{
					{
						ID:  "api-1",
						URL: "http://localhost:9001",
					},
					{
						ID:  "api-1",
						URL: "http://localhost:9002",
					},
				},
			},
			wantErr: true,
		},
		{
			name: "invalid upstream url",
			config: Config{
				Server: ServerConfig{
					Port: 8080,
				},
				LoadBalancer: LoadBalancerConfig{
					Strategy: RoundRobin,
				},
				Upstreams: []Upstream{
					{
						ID:  "api-1",
						URL: "invalid-url",
					},
				},
			},
			wantErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := Validate(&test.config)

			if test.wantErr && err == nil {
				t.Fatal("expected validation to fail")
			}

			if !test.wantErr && err != nil {
				t.Fatalf("expected validation to succeed: %v", err)
			}
		})
	}
}
