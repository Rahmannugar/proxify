package transport

import (
	"testing"

	"github.com/rahmannugar/proxify/internal/config"
)

func TestNew_ConfiguresConnectionPool(t *testing.T) {
	tests := []struct {
		name               string
		config             config.TransportConfig
		wantMaxIdle        int
		wantMaxIdlePerHost int
	}{
		{
			name:               "defaults",
			wantMaxIdle:        config.DefaultMaxIdleConns,
			wantMaxIdlePerHost: config.DefaultMaxIdleConnsPerHost,
		},
		{
			name: "configured values",
			config: config.TransportConfig{
				MaxIdleConns:        2000,
				MaxIdleConnsPerHost: 500,
			},
			wantMaxIdle:        2000,
			wantMaxIdlePerHost: 500,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			transport := New(test.config)

			if transport.MaxIdleConns != test.wantMaxIdle {
				t.Fatalf("expected MaxIdleConns %d, got %d", test.wantMaxIdle, transport.MaxIdleConns)
			}

			if transport.MaxIdleConnsPerHost != test.wantMaxIdlePerHost {
				t.Fatalf("expected MaxIdleConnsPerHost %d, got %d", test.wantMaxIdlePerHost, transport.MaxIdleConnsPerHost)
			}
		})
	}
}
