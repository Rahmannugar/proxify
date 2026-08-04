package registry

import (
	"testing"

	"github.com/rahmannugar/proxify/config"
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

	registry := New(upstreams)

	got := registry.List()

	if len(got) != len(upstreams) {
		t.Fatalf(
			"expected %d upstreams, got %d",
			len(upstreams),
			len(got),
		)
	}

	for i := range upstreams {
		if got[i] != upstreams[i] {
			t.Fatalf(
				"expected upstream %v, got %v",
				upstreams[i],
				got[i],
			)
		}
	}
}
