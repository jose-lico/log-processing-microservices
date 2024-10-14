package storelog

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/jose-lico/log-processing-microservices/common/logging"
	pb "github.com/jose-lico/log-processing-microservices/common/protos"

	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type Server struct {
	pb.UnimplementedLogServiceServer
	db *sql.DB
}

func NewServer(gRPC *grpc.Server, db *sql.DB) *Server {
	server := &Server{db: db}

	pb.RegisterLogServiceServer(gRPC, server)

	return server
}

func (s *Server) SubmitLog(ctx context.Context, in *pb.ProcessLogEntry) (*pb.LogResponse, error) {
	logging.Logger.Info("Received log entry", zap.Any("log", in))

	err := s.insertLogIntoDB(in)
	if err != nil {
		return &pb.LogResponse{Status: "500", Message: "Failed to write log entry"}, err
	}
	logging.Logger.Info("Stored log entry", zap.Any("log", in))

	return &pb.LogResponse{Status: "Success", Message: "Log entry stored"}, nil
}

func (s *Server) insertLogIntoDB(in *pb.ProcessLogEntry) error {
	additionalDataJSON, err := json.Marshal(in.AdditionalData)
	if err != nil {
		return fmt.Errorf("error marshaling additional data: %v", err)
	}

	_, err = s.db.Exec(`
		INSERT INTO process_log_entries (timestamp, level, message, user_id, additional_data, processed)
		VALUES ($1, $2, $3, $4, $5, $6)`,
		in.Timestamp, in.Level, in.Message, in.UserId, additionalDataJSON, in.Processed)

	if err != nil {
		return fmt.Errorf("error inserting log entry: %v", err)
	}

	return nil
}
