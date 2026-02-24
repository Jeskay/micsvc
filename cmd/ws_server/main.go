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
	"github.com/Jeskay/micsvc/internal/messager"
	transport "github.com/Jeskay/micsvc/internal/transport/websocket"
)

func main() {
	var cfg config.ServerConfig
	if err := godotenv.Load(); err != nil {
		log.Fatalf("failed to load env: %v", err)
	}
	if err := env.Parse(&cfg); err != nil {
		log.Fatal(err)
	}
	
	messageSvc := messager.NewMessagerSvc(func (s string) {
		log.Println(s)
	})
	server := transport.NewWebsocketServer(messageSvc)

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGTERM, syscall.SIGINT)

	done := make(chan struct{})

	go func() {
		<-c
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()
		err := server.Shutdown(ctx)
		if err != nil {
			log.Println(err)
		}
		close(done)
	}()

	if err := server.Run(cfg.Address()); err != nil {
		log.Fatal(err)
	}

	<-done
	log.Println("Server stopped")
}
