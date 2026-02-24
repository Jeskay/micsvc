package transport

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/Jeskay/micsvc/internal/transport/http/handlers"
	"github.com/Jeskay/micsvc/internal/user"
)

type HTTPServer struct {
	userSvc *user.Service
	server  *http.Server
}

func NewHTTPServer(userSvc *user.Service) *HTTPServer {
	return &HTTPServer{
		userSvc: userSvc,
	}
}

func (s *HTTPServer) Run(addr string) error {
	s.server = &http.Server{
		Addr:    addr,
		Handler: s.handler(),
	}
	return s.server.ListenAndServe()
}

func (s *HTTPServer) Shutdown(ctx context.Context) error {
	log.Println("Initiating server shutdown")
	if err := s.server.Shutdown(ctx); err != nil {
		return fmt.Errorf("server shutdown: %w", err)
	}
	return nil
}

func (s *HTTPServer) handler() http.Handler {
	userHandler := handlers.NewUserHandler(s.userSvc)

	m := http.NewServeMux()
	m.HandleFunc("POST /users", userHandler.Add())
	m.HandleFunc("PUT /users/{id}", userHandler.Update())
	m.HandleFunc("DELETE /users/{id}", userHandler.Delete())
	m.HandleFunc("GET /users", userHandler.GetAll())
	return m
}
