package transport

import (
	"net/http"

	"github.com/Jeskay/micsvc/internal/metrics"
	"github.com/Jeskay/micsvc/internal/transport/http/handlers"
	"github.com/Jeskay/micsvc/internal/transport/http/middleware"
	"github.com/Jeskay/micsvc/internal/user"
)

type UserHTTPOpts struct {
	Metrics *metrics.HTTPMetrics
}

func NewUserHTTP(userSvc *user.Service, addr string, opts UserHTTPOpts) *http.Server {
	userHandler := handlers.NewUserHandler(userSvc)

	m := http.NewServeMux()
	m.HandleFunc("POST /users", userHandler.Add())
	m.HandleFunc("PUT /users/{id}", userHandler.Update())
	m.HandleFunc("DELETE /users/{id}", userHandler.Delete())
	m.HandleFunc("GET /users", userHandler.GetAll())

	var handler http.Handler = m
	if opts.Metrics != nil {
		handler = middleware.Metrics(opts.Metrics, handler)
	}

	return &http.Server{
		Addr:    addr,
		Handler: handler,
	}
}

func NewMetricHTTP(metricSvc *metrics.Service, addr string) *http.Server {
	m := http.NewServeMux()
	m.Handle("/metrics", metricSvc.Handler())

	return &http.Server{
		Addr:    addr,
		Handler: m,
	}
}
