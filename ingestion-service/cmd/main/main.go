package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/jose-lico/log-processing-microservices/common/api"
	"github.com/jose-lico/log-processing-microservices/common/config"
	"github.com/jose-lico/log-processing-microservices/common/envs"
)

func main() {
	env := os.Getenv("ENV")

	if env == "LOCAL" {
		err := envs.LoadEnvs()
		if err != nil {
			log.Fatalf("[FATAL] Error loading .env: %v", err)
		}
	}

	cfg := config.NewRESTConfig()
	api := api.NewRESTServer(cfg)
	api.UseDefaultMiddleware()

	api.Router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello World!")
	})

	err := api.Run()
	if err != nil {
		log.Fatalf("[FATAL] Error launching API server: %v", err)
	}
}
