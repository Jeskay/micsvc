package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	env "github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"

	"github.com/Jeskay/micsvc/config"
	"github.com/Jeskay/micsvc/internal/auth"
	"github.com/Jeskay/micsvc/internal/db"
	transport "github.com/Jeskay/micsvc/internal/transport/grpc"
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
	authSvc := auth.NewAuthService(&cfg, memStore)

	server := transport.NewRPCServer(userSvc, authSvc)
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGTERM, syscall.SIGINT)

	done := make(chan struct{})

	go func() {
		<-c
		log.Println("Initiating server shutdown")
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()
		server.Shutdown(ctx)
		close(done)
	}()

	if err := server.Run(cfg.Address()); err != nil {
		log.Fatal(err)
	}

	<-done
	log.Println("Server stopped")
}
