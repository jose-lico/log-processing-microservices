package logservice

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

type server struct {
	pb.UnimplementedLogServiceServer
	db *sql.DB
}

func NewServer(gRPC *grpc.Server, db *sql.DB) {
	server := &server{db: db}

	pb.RegisterLogServiceServer(gRPC, server)
}

func (s *server) StoreLog(ctx context.Context, in *pb.StoreLogRequest) (*pb.StoreLogResponse, error) {
	logging.Logger.Info("Received log entry", zap.Any("log", in))

	err := s.insertLogIntoDB(in)
	if err != nil {
		return &pb.StoreLogResponse{Status: "500", Message: "Failed to write log entry"}, err
	}
	logging.Logger.Info("Stored log entry", zap.Any("log", in))

	return &pb.StoreLogResponse{Status: "Success", Message: "Log entry stored"}, nil
}

func (s *server) RetrieveLog(ctx context.Context, in *pb.RetrieveLogRequest) (*pb.RetrieveLogResponse, error) {
	logging.Logger.Info("Received logs request by id", zap.String("log", in.Id))

	logs, err := s.retrieveLogFromDB(in.Id, in.TimestampFrom, in.TimestampTo)

	if err != nil {
		return &pb.RetrieveLogResponse{Entries: nil, Status: nil}, nil
	}

	return &pb.RetrieveLogResponse{Entries: logs, Status: nil}, nil
}

func (s *server) insertLogIntoDB(in *pb.StoreLogRequest) error {
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

func (s *server) retrieveLogFromDB(id, from, to string) ([]*pb.StoreLogRequest, error) {
	var rows *sql.Rows
	var err error

	if from == "" {
		rows, err = s.db.Query("SELECT timestamp, level, message, user_id, additional_data, processed FROM process_log_entries WHERE user_id = $1", id)
	} else {
		rows, err = s.db.Query("SELECT timestamp, level, message, user_id, additional_data, processed FROM process_log_entries WHERE user_id = $1 AND timestamp BETWEEN $2 AND $3", id, from, to)
	}

	if err != nil {
		return nil, fmt.Errorf("error querying log entries: %v", err)
	}
	defer rows.Close()

	var entries []*pb.StoreLogRequest

	for rows.Next() {
		var entry pb.StoreLogRequest
		var additionalDataJSON []byte

		err := rows.Scan(&entry.Timestamp, &entry.Level, &entry.Message, &entry.UserId, &additionalDataJSON, &entry.Processed)
		if err != nil {
			return nil, fmt.Errorf("error scanning log entry: %v", err)
		}

		err = json.Unmarshal(additionalDataJSON, &entry.AdditionalData)
		if err != nil {
			return nil, fmt.Errorf("error unmarshaling additional data: %v", err)
		}

		entries = append(entries, &entry)
	}

	return entries, nil
}
