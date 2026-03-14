package middleware

import (
	"net/http"
	"time"

	"github.com/Jeskay/micsvc/internal/metrics"
)

func Metrics(httpMetric *metrics.HTTPMetrics, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t := time.Now()

		next.ServeHTTP(w, r)

		httpMetric.HTTPRequestLatency.WithLabelValues(r.Method).Observe(time.Since(t).Seconds())
	})
}
