package metrics

import "github.com/prometheus/client_golang/prometheus"

type HTTPMetrics struct {
	HTTPRequestLatency *prometheus.HistogramVec
}

func NewHTTPMetrics() *HTTPMetrics {
	return &HTTPMetrics{
		HTTPRequestLatency: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Name: "http_request_duration_seconds",
			Help: "Histogram of HTTP request latencies in seconds",
		}, []string{"method"}),
	}
}

func (m *HTTPMetrics) Register(r *prometheus.Registry) {
	r.MustRegister(m.HTTPRequestLatency)
}
