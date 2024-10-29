package processlog

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"github.com/jose-lico/log-processing-microservices/common/logging"
	pb "github.com/jose-lico/log-processing-microservices/common/protos"
	log_types "github.com/jose-lico/log-processing-microservices/common/types"

	"github.com/IBM/sarama"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type Consumer struct {
	wg     *sync.WaitGroup
	client pb.LogServiceClient
}

func NewConsumer(grpc *grpc.ClientConn, wg *sync.WaitGroup) *Consumer {
	consumer := &Consumer{wg: wg}

	consumer.client = pb.NewLogServiceClient(grpc)

	return consumer
}

func (c *Consumer) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

func (c *Consumer) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (c *Consumer) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		select {
		case <-sess.Context().Done():
			logging.Logger.Info("Session context canceled, exiting ConsumeClaim")
			return nil
		default:
		}

		logging.Logger.Info("Log Message received", zap.String("value", string(msg.Value)), zap.String("topic", msg.Topic), zap.Int32("partition", msg.Partition), zap.Int64("offset", msg.Offset))

		time.Sleep(5 * time.Second)

		var logEntry log_types.ProccessLogEntry

		err := json.Unmarshal([]byte(msg.Value), &logEntry)
		if err != nil {
			logging.Logger.Error("Invalid JSON format", zap.Error(err))
			return err
		}

		// Dummy business logic...
		logEntry.Processed = true

		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		entry := &pb.StoreLogRequest{
			Timestamp:      logEntry.Timestamp,
			Level:          logEntry.Level,
			Message:        logEntry.Message,
			UserId:         logEntry.UserID,
			AdditionalData: logEntry.AdditionalData,
			Processed:      logEntry.Processed,
		}

		// TODO: Handle status
		response, err := c.client.StoreLog(ctx, entry)
		if err != nil {
			logging.Logger.Error("Error submitting log", zap.Error(err))
			return err
		}

		logging.Logger.Info("Received storage response", zap.String("response", response.GetMessage()))

		sess.MarkMessage(msg, "")

		logging.Logger.Info("Log Message consumed", zap.String("value", string(msg.Value)), zap.String("topic", msg.Topic), zap.Int32("partition", msg.Partition), zap.Int64("offset", msg.Offset))
	}
	return nil
}
