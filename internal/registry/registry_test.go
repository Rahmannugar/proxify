package registry

import (
	"testing"

	"github.com/rahmannugar/proxify/internal/config"
)

func TestRegistry_List_ReturnsConfiguredUpstreams(t *testing.T) {
	upstreams := []config.Upstream{
		{
			ID:  "api-1",
			URL: "http://localhost:9001",
		},
		{
			ID:  "api-2",
			URL: "http://localhost:9002",
		},
	}

	registry, err := New(upstreams)
	if err != nil {
		t.Fatalf("expected registry creation to succeed: %v", err)
	}

	got := registry.List()

	if len(got) != len(upstreams) {
		t.Fatalf(
			"expected %d upstreams, got %d",
			len(upstreams),
			len(got),
		)
	}

	for i := range upstreams {
		if got[i].ID != upstreams[i].ID {
			t.Fatalf(
				"expected upstream ID %q, got %q",
				upstreams[i].ID,
				got[i].ID,
			)
		}

		if got[i].URL != upstreams[i].URL {
			t.Fatalf(
				"expected upstream URL %q, got %q",
				upstreams[i].URL,
				got[i].URL,
			)
		}
	}
}
