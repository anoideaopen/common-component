package baseprometheus

import (
	"context"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// MetricsHandler returns a handler for prometheus metrics
func MetricsHandler(_ context.Context) http.Handler {
	return promhttp.Handler()
}
