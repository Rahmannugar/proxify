package loadbalancer

import (
	"testing"

	"github.com/rahmannugar/proxify/internal/config"
	"github.com/rahmannugar/proxify/internal/registry"
)

func TestLeastOutstanding_Next_ReturnsLeastBusyUpstream(t *testing.T) {
	registry, err := registry.New([]config.Upstream{
		{
			ID:  "api-1",
			URL: "http://localhost:9001",
		},
		{
			ID:  "api-2",
			URL: "http://localhost:9002",
		},
		{
			ID:  "api-3",
			URL: "http://localhost:9003",
		},
	})
	if err != nil {
		t.Fatalf("expected registry creation to succeed: %v", err)
	}

	upstreams := registry.List()

	/*
		api-1 -> 3 active requests
		api-2 -> 1 active request
		api-3 -> 0 active requests
	*/

	upstreams[0].IncrementActiveRequests()
	upstreams[0].IncrementActiveRequests()
	upstreams[0].IncrementActiveRequests()

	upstreams[1].IncrementActiveRequests()

	lb := NewLeastOutstanding(registry)

	upstream, err := lb.Next()
	if err != nil {
		t.Fatalf("expected upstream, got error: %v", err)
	}

	if upstream.ID != "api-3" {
		t.Fatalf(
			"expected api-3, got %q",
			upstream.ID,
		)
	}
}

func TestLeastOutstanding_Next_SkipsUnhealthyUpstreams(t *testing.T) {
	registry, err := registry.New([]config.Upstream{
		{
			ID:  "api-1",
			URL: "http://localhost:9001",
		},
		{
			ID:  "api-2",
			URL: "http://localhost:9002",
		},
	})
	if err != nil {
		t.Fatalf("expected registry creation to succeed: %v", err)
	}

	upstreams := registry.List()

	upstreams[0].SetHealthy(false)

	lb := NewLeastOutstanding(registry)

	upstream, err := lb.Next()
	if err != nil {
		t.Fatalf("expected upstream, got error: %v", err)
	}

	if upstream.ID != "api-2" {
		t.Fatalf(
			"expected api-2, got %q",
			upstream.ID,
		)
	}
}

func TestLeastOutstanding_Next_ReturnsErrorWhenNoHealthyUpstreams(t *testing.T) {
	registry, err := registry.New([]config.Upstream{
		{
			ID:  "api-1",
			URL: "http://localhost:9001",
		},
		{
			ID:  "api-2",
			URL: "http://localhost:9002",
		},
	})
	if err != nil {
		t.Fatalf("expected registry creation to succeed: %v", err)
	}

	for _, upstream := range registry.List() {
		upstream.SetHealthy(false)
	}

	lb := NewLeastOutstanding(registry)

	_, err = lb.Next()
	if err == nil {
		t.Fatal("expected an error when no healthy upstreams are available")
	}
}
