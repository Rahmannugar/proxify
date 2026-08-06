package exporter

import (
	"encoding/json"
	"net/http"

	"github.com/rahmannugar/proxify/internal/metrics"
)

// JSONMetrics exports operational metrics as JSON.
type JSONMetrics struct {
	metrics *metrics.Metrics
}

// NewJSONMetrics creates a new JSON metrics exporter.
func NewJSONMetrics(metrics *metrics.Metrics) *JSONMetrics {
	return &JSONMetrics{
		metrics: metrics,
	}
}

// ServeHTTP exposes the current metrics snapshot.
func (e *JSONMetrics) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(e.metrics.Snapshot()); err != nil {
		http.Error(
			w,
			"failed to encode metrics",
			http.StatusInternalServerError,
		)

		return
	}
}