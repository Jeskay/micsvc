package transport

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Jeskay/micsvc/internal/messager"
)

type WebsocketServer struct {
	msgSvc *messageServer
	server *http.Server
}

func NewWebsocketServer(msgSvc *messager.Service) *WebsocketServer {
	return &WebsocketServer{msgSvc: NewMessageServer(msgSvc)}
}

func (s *WebsocketServer) Run(addr string) error {
	s.server = &http.Server{
		Handler: s.msgSvc,
		Addr:    addr,
	}
	return s.server.ListenAndServe()
}

func (s *WebsocketServer) Shutdown(ctx context.Context) error {
	if err := s.server.Shutdown(ctx); err != nil {
		return fmt.Errorf("server shutdown: %w", err)
	}
	return nil
}
