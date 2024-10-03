package main

import (
	"os"

	"github.com/jose-lico/log-processing-microservices/ingestion-service/ingest_log_service"

	"github.com/jose-lico/log-processing-microservices/common/api"
	"github.com/jose-lico/log-processing-microservices/common/config"
	"github.com/jose-lico/log-processing-microservices/common/envs"
	"github.com/jose-lico/log-processing-microservices/common/logging"
	"github.com/jose-lico/log-processing-microservices/common/middleware"
	"go.uber.org/zap"

	chi_middleware "github.com/go-chi/chi/v5/middleware"
)

func main() {
	env := os.Getenv("ENV")

	if env == "LOCAL" {
		err := envs.LoadEnvs()
		if err != nil {
			panic(err)
		}
	}

	cfg := config.NewRESTConfig()
	api := api.NewRESTServer(cfg)
	api.Router.Use(middleware.LoggingMiddleware())
	api.Router.Use(chi_middleware.Recoverer)

	api.Router.Post("/", ingest_log_service.IngestLog)

	err := api.Run()
	if err != nil {
		logging.Logger.Fatal("Error launching HTTP Server", zap.Error(err))
		panic(err)
	}
}
