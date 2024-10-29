package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/jose-lico/log-processing-microservices/ingestion-service/services/ingestlog"

	"github.com/jose-lico/log-processing-microservices/common/api"
	"github.com/jose-lico/log-processing-microservices/common/config"
	"github.com/jose-lico/log-processing-microservices/common/envs"
	"github.com/jose-lico/log-processing-microservices/common/kafka"
	"github.com/jose-lico/log-processing-microservices/common/logging"
	"github.com/jose-lico/log-processing-microservices/common/middleware"

	chi_middleware "github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
)

func main() {
	env := os.Getenv("ENV")

	logging.CreateLogger(env)
	defer logging.Logger.Sync()

	if env == "LOCAL" {
		err := envs.LoadEnvs()
		if err != nil {
			logging.Logger.Fatal("Failed to load envs", zap.Error(err))
		}
	}

	ctx, cancel := context.WithCancel(context.Background())
	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigterm
		logging.Logger.Info("Received termination signal, shutting down...")
		cancel()
	}()

	kafkaHost := os.Getenv("KAFKA_HOST")
	kafkaPort := os.Getenv("KAFKA_PORT")
	if kafkaHost == "" || kafkaPort == "" {
		logging.Logger.Fatal("KAFKA_HOST and KAFKA_PORT must be set",
			zap.String("KAFKA_HOST", kafkaHost),
			zap.String("KAFKA_PORT", kafkaPort))
	}
	producer, err := kafka.NewAsyncProducer(ctx, "ingestion-service", []string{fmt.Sprintf("%s:%s", kafkaHost, kafkaPort)})
	if err != nil {
		logging.Logger.Fatal("Failed to create Kafka producer", zap.Error(err))
	}

	cfg := config.NewRESTConfig()
	api := api.NewRESTServer()
	api.Router.Use(middleware.LoggingMiddleware())
	api.Router.Use(chi_middleware.Recoverer)

	ingestLogService := ingestlog.NewService(producer)
	ingestLogService.RegisterRoutes(api.Router)

	err = api.Run(ctx, cfg)
	if err != nil && !errors.Is(err, context.Canceled) && !errors.Is(err, context.DeadlineExceeded) {
		logging.Logger.Error("Error running REST server", zap.Error(err))
	}

	producer.Close()

	logging.Logger.Info("Ingestion Service has shutdown")
}
