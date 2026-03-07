package main

import (
	"context"
	"fmt"
	"log"
	"os/signal"
	"syscall"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"

	"github.com/Jeskay/micsvc/config"
	"github.com/Jeskay/micsvc/internal/messager"
	transport "github.com/Jeskay/micsvc/internal/transport/websocket"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}
	var cfg config.ClientConfig
	if err := env.Parse(&cfg); err != nil {
		log.Fatal(err)
	}
	addr := fmt.Sprintf("ws://%s:%s", cfg.Host, cfg.Port)
	c := transport.NewWebSocketClient("Reader")
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	if err := c.Connect(ctx, addr); err != nil {
		log.Fatalf("failed to connect: %v", err)
	}

	go func() {
		startReading(ctx, c.Out)
		stop()
	}()
	<-ctx.Done()
	log.Println("finished")
}

func startReading(ctx context.Context, out chan messager.Message) {
	for {
		select {
		case <-ctx.Done():
			return
		case m, ok := <-out:
			if !ok {
				return
			}
			if !m.Binary {
				log.Printf("%s: %s", m.AuthorID, string(m.Data))
			}
		}
	}
}
