package main

import (
	"context"
	"errors"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/IBM/sarama"
	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
	"golang.org/x/sync/errgroup"

	"github.com/Jeskay/micsvc/config"
	"github.com/Jeskay/micsvc/internal/db"
	"github.com/Jeskay/micsvc/internal/metrics"
	broker "github.com/Jeskay/micsvc/internal/transport/broker/producer"
	transport "github.com/Jeskay/micsvc/internal/transport/http"
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
				return slog.Attr{Key: "log.level", Value: slog.StringValue(a.Value.String())}
			}
			if a.Key == slog.MessageKey {
				return slog.Attr{Key: "message", Value: a.Value}
			}
			return a
		},
	}
	logger := slog.New(slog.NewJSONHandler(os.Stdout, opts))
	if err := godotenv.Load(); err != nil && !os.IsNotExist(err) {
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

	HTTPcollector := metrics.NewHTTPMetrics()
	userSvc := user.NewUserService(logger.WithGroup("user-service"), memStore, eProducer)
	metricSvc := metrics.NewService(HTTPcollector)

	usrServer := transport.NewUserHTTP(
		userSvc,
		cfg.Address(),
		transport.UserHTTPOpts{Metrics: HTTPcollector},
	)
	metricServer := transport.NewMetricHTTP(metricSvc, cfg.MetricAddress())

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer stop()
	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		return metricServer.ListenAndServe()
	})

	g.Go(func() error {
		return usrServer.ListenAndServe()
	})

	g.Go(func() error {
		return eProducer.Run()
	})

	g.Go(func() error {
		<-ctx.Done()
		eProducer.Shutdown()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()
		return errors.Join(usrServer.Shutdown(shutdownCtx), metricServer.Shutdown(shutdownCtx))
	})

	if err := g.Wait(); err != nil {
		log.Printf("Execution interrupted by: %v", err)
	}

	log.Println("Server stopped")
}
