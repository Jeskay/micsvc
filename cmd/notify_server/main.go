package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/IBM/sarama"
	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
	"golang.org/x/sync/errgroup"

	"github.com/Jeskay/micsvc/config"
	"github.com/Jeskay/micsvc/internal/messager"
	"github.com/Jeskay/micsvc/internal/notify"
	broker "github.com/Jeskay/micsvc/internal/transport/broker/consumer"
	transport "github.com/Jeskay/micsvc/internal/transport/websocket"
)

func main() {
	var cfg config.ServerConfig
	if err := godotenv.Load(); err != nil && !os.IsNotExist(err) {
		log.Fatal(err)
	}
	if err := env.Parse(&cfg); err != nil {
		log.Fatal(err)
	}
	msgSvc := messager.NewMessagerSvc(func(s string) { log.Println(s) }, cfg.ConnectionTimeout())
	consumer, err := sarama.NewConsumer([]string{cfg.KafkaAddress}, sarama.NewConfig())
	if err != nil {
		log.Fatalf("failed to init consumer: %v", err)
	}
	eConsumer, err := broker.NewEventConsumer(consumer, cfg.EventTopic, 0, sarama.OffsetOldest)
	if err != nil {
		log.Fatalf("failed to start consumer: %v", err)
	}
	notifySvc := notify.NewService(msgSvc, eConsumer)
	server := transport.NewWebsocketServer(msgSvc)
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer stop()
	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		return notifySvc.Run("Notification")
	})

	g.Go(func() error {
		return server.Run(cfg.Address())
	})

	g.Go(func() error {
		<-ctx.Done()
		notifySvc.Close()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()
		return server.Shutdown(shutdownCtx)
	})

	if err := g.Wait(); err != nil {
		log.Printf("Exiting with error: %v", err)
	}
	log.Println("Server stopped")
}
