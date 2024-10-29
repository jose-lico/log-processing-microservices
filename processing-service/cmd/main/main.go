package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/jose-lico/log-processing-microservices/processing-service/consumers/processlog"

	"github.com/jose-lico/log-processing-microservices/common/envs"
	"github.com/jose-lico/log-processing-microservices/common/kafka"
	"github.com/jose-lico/log-processing-microservices/common/logging"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	env := os.Getenv("ENV")

	logging.CreateLogger(env)
	defer logging.Logger.Sync()

	if env == "LOCAL" {
		err := envs.LoadEnvs()
		if err != nil {
			panic(err)
		}
	}

	storageHost := os.Getenv("STORAGE_HOST")
	storagePort := os.Getenv("STORAGE_PORT")
	kafkaHost := os.Getenv("KAFKA_HOST")
	kafkaPort := os.Getenv("KAFKA_PORT")

	if storageHost == "" || storagePort == "" || kafkaHost == "" || kafkaPort == "" {
		logging.Logger.Fatal("STORAGE_HOST, STORAGE_PORT and KAFKA_HOST, KAFKA_PORT must be set",
			zap.String("STORAGE_HOST", storageHost), zap.String("STORAGE_PORT", storagePort),
			zap.String("KAFKA_HOST", kafkaHost), zap.String("KAFKA_PORT", kafkaPort))
	}

	conn, err := grpc.NewClient(fmt.Sprintf("%s:%s", storageHost, storagePort), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logging.Logger.Fatal("Could not create gRPC client", zap.Error(err))
	}
	logging.Logger.Info("Created gRPC Client")

	consumerGroup, err := kafka.CreateKafkaConsumer([]string{fmt.Sprintf("%s:%s", kafkaHost, kafkaPort)}, "processing-service", "logs_group")
	if err != nil {
		logging.Logger.Fatal("Failed to start Kafka producer", zap.Error(err))
	}

	ctx, cancel := context.WithCancel(context.Background())
	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigterm
		logging.Logger.Info("Received termination signal, shutting down...")
		cancel()
	}()

	// var consumerWg sync.WaitGroup
	consumer := processlog.NewConsumer(conn, nil)
	// consumerWg.Add(1)

	if err := consumerGroup.Consume(ctx, []string{"logs"}, consumer); err != nil {
		logging.Logger.Error("Error from consumer", zap.Error(err))
	}

	if err := consumerGroup.Close(); err != nil {
		logging.Logger.Error("Error closing consumer group", zap.Error(err))
	}

	logging.Logger.Info("Closing gRPC Client")
	if err := conn.Close(); err != nil {
		logging.Logger.Error("Error closing gRPC connection", zap.Error(err))
	}

	logging.Logger.Info("Processing service has shut down")
}
