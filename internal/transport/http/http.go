package transport

import (
	"net/http"

	"github.com/Jeskay/micsvc/internal/transport/http/handlers"
	"github.com/Jeskay/micsvc/internal/user"
)

func handler(svc *user.Service) http.Handler {
	userHandler := handlers.NewUserHandler(svc)

	m := http.NewServeMux()
	m.HandleFunc("POST /users", userHandler.Add())
	m.HandleFunc("PUT /users/{id}", userHandler.Update())
	m.HandleFunc("DELETE /users/{id}", userHandler.Delete())
	m.HandleFunc("GET /users", userHandler.GetAll())
	return m
}

func RunHTTP(addr string, svc *user.Service, onClose func(error)) {
	if err := http.ListenAndServe(addr, handler(svc)); err != nil {
		onClose(err)
	}
}
