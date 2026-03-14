package metrics

import "github.com/prometheus/client_golang/prometheus"

type GRPCMetrics struct {
	GRPCReuqestCounter *prometheus.CounterVec
	GRPCRequestLatency *prometheus.HistogramVec
}

func NewGRPCMetrics() *GRPCMetrics {
	return &GRPCMetrics{
		GRPCReuqestCounter: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: "grpc_requests_total",
			Help: "Total amount of gRPC requests received.",
		}, []string{"method", "status"}),
		GRPCRequestLatency: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Name: "grpc_request_latency_seconds",
			Help: "Histogram of gRPC request latencies in seconds.",
		}, []string{"method"}),
	}
}

func (m *GRPCMetrics) Register(r *prometheus.Registry) {
	r.MustRegister(m.GRPCRequestLatency, m.GRPCReuqestCounter)
}
