package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"

	"github.com/Jeskay/micsvc/config"
	"github.com/Jeskay/micsvc/internal/db"
	transport "github.com/Jeskay/micsvc/internal/transport/http"
	"github.com/Jeskay/micsvc/internal/user"
)

func main() {
	var cfg config.ServerConfig
	if err := godotenv.Load(); err != nil {
		log.Fatalf("failed to load env: %v", err)
	}
	err := env.Parse(&cfg)
	if err != nil {
		log.Fatalf("failed to parse params: %v", err)
	}
	memStore := db.NewUserStorage()
	userSvc := user.NewUserService(memStore)
	server := transport.NewHTTPServer(userSvc)
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGTERM, syscall.SIGINT)

	done := make(chan struct{})

	go func() {
		<-c
		log.Println("Initiating server shutdown")
		ctx, cancel := context.WithTimeout(context.Background(), time.Second * 3)
		defer cancel()
		if err:= server.Shutdown(ctx); err != nil {
			log.Fatalf("Failed to gracefully shutdown: %v",err)
		}
		close(done)
	}()

	if err := server.Run(cfg.Address()); err != nil {
		log.Fatal(err)
	}

	<-done
	log.Println("Server stopped")
}
