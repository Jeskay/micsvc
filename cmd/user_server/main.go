package main

import (
	"context"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/IBM/sarama"
	env "github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
	"golang.org/x/sync/errgroup"

	"github.com/Jeskay/micsvc/config"
	"github.com/Jeskay/micsvc/internal/auth"
	"github.com/Jeskay/micsvc/internal/db"
	broker "github.com/Jeskay/micsvc/internal/transport/broker/producer"
	transport "github.com/Jeskay/micsvc/internal/transport/grpc"
	"github.com/Jeskay/micsvc/internal/user"
)

func main() {
	var cfg config.ServerConfig
	opts := &slog.HandlerOptions{
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey {
				return slog.Attr{Key: "@timestamp", Value: a.Value}
			}
			if a.Key == slog.LevelKey {
				return slog.Attr{Key: "level", Value: slog.StringValue(a.Value.String())}
			}
			return a
		},
	}
	logger := slog.New(slog.NewJSONHandler(os.Stdout, opts))
	if err := godotenv.Load(); err != nil {
		log.Fatalf("failed to load env: %v", err)
	}
	err := env.Parse(&cfg)
	if err != nil {
		log.Fatalf("failed to parse params: %v", err)
	}

	memStore := db.NewUserStorage()
	producer, err := sarama.NewAsyncProducer([]string{cfg.KafkaAddress}, sarama.NewConfig())
	if err != nil {
		log.Fatalf("failed to set up event producer: %v", err)
	}
	eProducer := broker.NewEventProducer(producer, cfg.EventTopic, cfg.ConnectionTimeout())
	userSvc := user.NewUserService(logger, memStore, eProducer)
	authSvc := auth.NewAuthService(&cfg, memStore)

	server := transport.NewRPCServer(userSvc, authSvc)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		return server.Run(cfg.Address())
	})

	g.Go(func() error {
		return eProducer.Run()
	})

	g.Go(func() error {
		<-ctx.Done()
		eProducer.Shutdown()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()
		server.Shutdown(shutdownCtx)
		return nil
	})

	if err := g.Wait(); err != nil {
		log.Printf("Execution interrupted: %v", err)
	}

	log.Println("Server stopped")
}
