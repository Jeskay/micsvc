package metrics

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type MetricCollector interface {
	Register(r *prometheus.Registry)
}

type Service struct {
	Registry *prometheus.Registry
}

func NewService(collectors ...MetricCollector) *Service {
	registry := prometheus.NewRegistry()
	m := &Service{
		Registry: registry,
	}
	m.Register(collectors...)
	m.RegisterDefault()
	return m
}

func (s *Service) Register(collectors ...MetricCollector) {
	for _, c := range collectors {
		c.Register(s.Registry)
	}
}

func (s *Service) RegisterDefault() {
	s.Registry.MustRegister(
		collectors.NewGoCollector(),
		collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
	)
}

func (s *Service) Handler() http.Handler {
	return promhttp.HandlerFor(s.Registry, promhttp.HandlerOpts{})
}
