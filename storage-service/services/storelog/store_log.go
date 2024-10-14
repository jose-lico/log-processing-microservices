package storelog

import (
	"context"

	"github.com/jose-lico/log-processing-microservices/common/logging"
	pb "github.com/jose-lico/log-processing-microservices/common/protos"
	"go.uber.org/zap"

	"google.golang.org/grpc"
)

type Server struct {
	pb.UnimplementedLogServiceServer
}

func NewServer(gRPC *grpc.Server) *Server {
	server := &Server{}

	pb.RegisterLogServiceServer(gRPC, server)

	return &Server{}
}

func (s *Server) SubmitLog(ctx context.Context, in *pb.ProcessLogEntry) (*pb.LogResponse, error) {
	logging.Logger.Info("Received log entry", zap.Any("log", in))

	// err := insertLogIntoDB(in)
	// if err != nil {
	//   return &pb.LogResponse{Status: "Error", Message: "Failed to write log entry"}, err
	// }

	return &pb.LogResponse{Status: "Success", Message: "Log entry stored"}, nil
}
