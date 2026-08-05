package health

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/rahmannugar/proxify/config"
	"github.com/rahmannugar/proxify/registry"
)

func TestChecker_CheckUpstream_Healthy(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}),
	)
	defer server.Close()

	registry, err := registry.New([]config.Upstream{
		{
			ID:  "api-1",
			URL: server.URL,
		},
	})
	if err != nil {
		t.Fatalf("expected registry creation to succeed: %v", err)
	}

	checker := New(
		registry,
		"/",
		time.Second,
		time.Second,
		http.StatusOK,
	)

	upstream := registry.List()[0]

	checker.checkUpstream(upstream)

	if !upstream.IsHealthy() {
		t.Fatal("expected upstream to be healthy")
	}
}

func TestChecker_CheckUpstream_Unreachable(t *testing.T) {
	registry, err := registry.New([]config.Upstream{
		{
			ID:  "api-1",
			URL: "http://127.0.0.1:65535",
		},
	})
	if err != nil {
		t.Fatalf("expected registry creation to succeed: %v", err)
	}

	checker := New(
		registry,
		"/health",
		time.Second,
		time.Second,
		http.StatusOK,
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

	registry, err := registry.New([]config.Upstream{
		{
			ID:  "api-1",
			URL: server.URL,
		},
	})
	if err != nil {
		t.Fatalf("expected registry creation to succeed: %v", err)
	}

	checker := New(
		registry,
		"/",
		time.Second,
		time.Second,
		http.StatusOK,
	)

	upstream := registry.List()[0]

	checker.checkUpstream(upstream)

	if upstream.IsHealthy() {
		t.Fatal("expected upstream to be unhealthy")
	}
}