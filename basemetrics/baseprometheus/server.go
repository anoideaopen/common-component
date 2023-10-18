package baseprometheus

import (
	"context"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func MetricsHandler(_ context.Context) http.Handler {
	return promhttp.Handler()
}
