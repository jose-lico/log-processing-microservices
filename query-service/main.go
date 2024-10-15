package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/jose-lico/log-processing-microservices/common/api"
	"github.com/jose-lico/log-processing-microservices/common/config"
	"github.com/jose-lico/log-processing-microservices/common/envs"
	"github.com/jose-lico/log-processing-microservices/common/logging"
	"github.com/jose-lico/log-processing-microservices/common/middleware"
	pb "github.com/jose-lico/log-processing-microservices/common/protos"

	chi_middleware "github.com/go-chi/chi/v5/middleware"
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
	// from := r.PathValue("from")
	// to := r.PathValue("to")

	err := isValidUUID4(id)
	if err != nil {
		api.WriteJSON(w, http.StatusBadRequest, map[string]interface{}{
			"status":  "error",
			"message": "ID not valid",
			"error":   err,
		})
		return
	}

	// err := isValidTimestamps(from, to)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	request := &pb.RetrieveLogRequest{
		Id: id,
	}

	response, _ := client.RetrieveLogByID(ctx, request)

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

// func isValidTimestamps(from, to string) error {
// 	return nil
// }
