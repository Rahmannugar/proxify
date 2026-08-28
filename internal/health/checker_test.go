package health

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/rahmannugar/proxify/internal/config"
	"github.com/rahmannugar/proxify/internal/metrics"
	"github.com/rahmannugar/proxify/internal/registry"
	"github.com/rahmannugar/proxify/internal/testutil"
	"github.com/rahmannugar/proxify/internal/transport"
)

func TestChecker_CheckUpstream_Healthy(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}),
	)
	defer server.Close()

	registry, err := registry.New(
		[]config.Upstream{
			{
				ID:  "api-1",
				URL: server.URL,
			},
		},
		testutil.CircuitBreakerConfig(),
		transport.New(config.TransportConfig{}),
	)
	if err != nil {
		t.Fatalf("expected registry creation to succeed: %v", err)
	}

	checker := New(
		registry,
		metrics.New(),
		"/",
		time.Second,
		time.Second,
		http.StatusOK,
		3,
	)

	upstream := registry.List()[0]

	checker.checkUpstream(upstream)

	if !upstream.IsHealthy() {
		t.Fatal("expected upstream to be healthy")
	}
}

func TestChecker_CheckUpstream_Unreachable(t *testing.T) {
	registry, err := registry.New(
		[]config.Upstream{
			{
				ID:  "api-1",
				URL: "http://127.0.0.1:65535",
			},
		},
		testutil.CircuitBreakerConfig(),
		transport.New(config.TransportConfig{}),
	)
	if err != nil {
		t.Fatalf("expected registry creation to succeed: %v", err)
	}

	checker := New(
		registry,
		metrics.New(),
		"/health",
		time.Second,
		time.Second,
		http.StatusOK,
		1,
	)

	upstream := registry.List()[0]

	checker.checkUpstream(upstream)

	if upstream.IsHealthy() {
		t.Fatal("expected upstream to be unhealthy")
	}
}

func TestChecker_CheckUpstream_UnexpectedStatus(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}),
	)
	defer server.Close()

	registry, err := registry.New(
		[]config.Upstream{
			{
				ID:  "api-1",
				URL: server.URL,
			},
		},
		testutil.CircuitBreakerConfig(),
		transport.New(config.TransportConfig{}),
	)
	if err != nil {
		t.Fatalf("expected registry creation to succeed: %v", err)
	}

	checker := New(
		registry,
		metrics.New(),
		"/",
		time.Second,
		time.Second,
		http.StatusOK,
		1,
	)

	upstream := registry.List()[0]

	checker.checkUpstream(upstream)

	if upstream.IsHealthy() {
		t.Fatal("expected upstream to be unhealthy")
	}
}

func TestChecker_CheckUpstream_RequiresConsecutiveFailures(t *testing.T) {
	status := http.StatusInternalServerError
	server := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(status)
		}),
	)
	defer server.Close()

	registry, err := registry.New(
		[]config.Upstream{{ID: "api-1", URL: server.URL}},
		testutil.CircuitBreakerConfig(),
		transport.New(config.TransportConfig{}),
	)
	if err != nil {
		t.Fatalf("expected registry creation to succeed: %v", err)
	}

	checker := New(
		registry,
		metrics.New(),
		"/",
		time.Second,
		time.Second,
		http.StatusOK,
		3,
	)
	upstream := registry.List()[0]

	checker.checkUpstream(upstream)
	checker.checkUpstream(upstream)

	if !upstream.IsHealthy() {
		t.Fatal("expected upstream to remain healthy below retry threshold")
	}

	status = http.StatusOK
	checker.checkUpstream(upstream)
	status = http.StatusInternalServerError
	checker.checkUpstream(upstream)

	if !upstream.IsHealthy() {
		t.Fatal("expected successful check to reset consecutive failures")
	}

	checker.checkUpstream(upstream)
	checker.checkUpstream(upstream)

	if upstream.IsHealthy() {
		t.Fatal("expected upstream to become unhealthy at retry threshold")
	}
}
