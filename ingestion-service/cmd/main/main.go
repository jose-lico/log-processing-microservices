package main

import (
	"fmt"
	"os"

	"github.com/jose-lico/log-processing-microservices/ingestion-service/ingest_log"

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

	if env == "LOCAL" {
		err := envs.LoadEnvs()
		if err != nil {
			panic(err)
		}
	}

	logging.CreateLogger()
	defer logging.Logger.Sync()

	kafkaHost := os.Getenv("KAFKA_HOST")
	kafkaPort := os.Getenv("KAFKA_PORT")
	producer, err := kafka.CreateKafkaProducer([]string{fmt.Sprintf("%s:%s", kafkaHost, kafkaPort)})
	if err != nil {
		logging.Logger.Fatal("Failed to start Kafka producer", zap.Error(err))
	}
	defer producer.Close()

	cfg := config.NewRESTConfig()
	api := api.NewRESTServer(cfg)
	api.Router.Use(middleware.LoggingMiddleware())
	api.Router.Use(chi_middleware.Recoverer)

	ingestLogService := ingest_log.NewService(producer)
	ingestLogService.RegisterRoutes(api.Router)

	err = api.Run()
	if err != nil {
		logging.Logger.Fatal("Error launching HTTP Server", zap.Error(err))
	}
}
