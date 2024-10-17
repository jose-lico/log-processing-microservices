package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/jose-lico/log-processing-microservices/common/api"
	"github.com/jose-lico/log-processing-microservices/common/config"
	"github.com/jose-lico/log-processing-microservices/common/envs"
	"github.com/jose-lico/log-processing-microservices/common/logging"
	"github.com/jose-lico/log-processing-microservices/common/middleware"
	pb "github.com/jose-lico/log-processing-microservices/common/protos"

	chi_middleware "github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var client pb.LogServiceClient

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

	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	address := fmt.Sprintf("%s:%s", os.Getenv("STORAGE_HOST"), os.Getenv("STORAGE_PORT"))
	conn, err := grpc.NewClient(address, opts...)
	if err != nil {
		logging.Logger.Fatal("Could not create gRPC client", zap.Error(err))
	}
	defer conn.Close()
	logging.Logger.Info("Created gRPC Client")

	client = pb.NewLogServiceClient(conn)

	cfg := config.NewRESTConfig()
	api := api.NewRESTServer(cfg)
	api.Router.Use(middleware.LoggingMiddleware())
	api.Router.Use(chi_middleware.Recoverer)

	api.Router.Get("/", getLogs)

	err = api.Run()
	if err != nil {
		logging.Logger.Fatal("Error launching HTTP Server", zap.Error(err))
	}
}

func getLogs(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	from := r.URL.Query().Get("from")
	to := r.URL.Query().Get("to")

	err := isValidUUID4(id)
	if err != nil {
		api.WriteJSON(w, http.StatusBadRequest, map[string]interface{}{
			"status":  "error",
			"message": "ID not valid",
			"error":   err.Error(),
		})
		return
	}

	request := &pb.RetrieveLogRequest{
		Id: id,
	}

	if from != "" || to != "" {
		err = isValidTimestamps(from, to)
		if err != nil {
			api.WriteJSON(w, http.StatusBadRequest, map[string]interface{}{
				"status":  "error",
				"message": "Timestamps not valid",
				"error":   err.Error(),
			})
			return
		}

		request.TimestampFrom = from
		request.TimestampTo = to
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	response, err := client.RetrieveLog(ctx, request)
	if err != nil {
		api.WriteJSON(w, http.StatusBadRequest, map[string]interface{}{
			"status":  "error",
			"message": "could not retrieve logs from storage service",
			"error":   err.Error(),
		})
		return
	}

	api.WriteJSON(w, http.StatusBadRequest, map[string]interface{}{
		"status": "ok",
		"logs":   response.Entries,
	})
}

func isValidUUID4(id string) error {
	if id == "" {
		return errors.New("id is empty")
	}

	u, err := uuid.Parse(id)
	if err != nil {
		return err
	}

	if u.Version() == 4 {
		return nil
	} else {
		return errors.New("id is not uuid v4")
	}
}

func isValidTimestamps(from, to string) error {
	if from == "" && to != "" {
		return errors.New("missing from date")
	} else if from != "" && to == "" {
		return errors.New("missing to date")
	}

	layout := time.RFC3339

	fromTime, err := time.Parse(layout, from)
	if err != nil {
		return fmt.Errorf("error parsing from: %v", err)
	}

	toTime, err := time.Parse(layout, to)
	if err != nil {
		return fmt.Errorf("error parsing to: %v", err)
	}

	if toTime.Before(fromTime) || toTime.Equal(fromTime) {
		return errors.New("to time is after or same as from time, invalid timeframe")
	}

	return nil
}
