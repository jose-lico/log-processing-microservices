package main

import (
	"net"
	"os"

	"github.com/jose-lico/log-processing-microservices/storage-service/database"
	"github.com/jose-lico/log-processing-microservices/storage-service/services/logservice"

	"github.com/jose-lico/log-processing-microservices/common/envs"
	"github.com/jose-lico/log-processing-microservices/common/logging"

	"go.uber.org/zap"
	"google.golang.org/grpc"
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

	db, err := database.NewPostresDB()
	if err != nil {
		logging.Logger.Fatal("Failed to connect to Postgres", zap.Error(err))
	}
	defer db.Close()

	logging.Logger.Info("Connected to Postgres")

	port := ":" + os.Getenv("PORT")

	lis, err := net.Listen("tcp", port)
	if err != nil {
		logging.Logger.Fatal("Failed to listen", zap.Error(err))
	}

	opts := []grpc.ServerOption{}
	grpcServer := grpc.NewServer(opts...)

	logservice.NewServer(grpcServer, db)

	logging.Logger.Info("gRPC Server is listening", zap.String("Port", port))

	if err := grpcServer.Serve(lis); err != nil {
		logging.Logger.Fatal("Failed to serve", zap.Error(err))
	}
}
