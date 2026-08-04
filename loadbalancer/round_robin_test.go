package loadbalancer

import (
	"sync"
	"testing"

	"github.com/rahmannugar/proxify/config"
	"github.com/rahmannugar/proxify/registry"
)

func TestRoundRobin_Next_CyclesThroughUpstreams(t *testing.T) {
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
	lb := NewRoundRobin(registry)

	expected := []string{
		"api-1",
		"api-2",
		"api-3",
		"api-1",
		"api-2",
		"api-3",
	}

	for _, id := range expected {
		upstream := lb.Next()

		if upstream.ID != id {
			t.Fatalf(
				"expected upstream %q, got %q",
				id,
				upstream.ID,
			)
		}
	}
}

func TestRoundRobin_Next_ConcurrentAccess(t *testing.T) {
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
	lb := NewRoundRobin(registry)

	const requests = 100

	var wg sync.WaitGroup
	wg.Add(requests)

	for range requests {
		go func() {
			defer wg.Done()

			_ = lb.Next()
		}()
	}

	wg.Wait()
}
