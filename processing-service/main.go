package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/jose-lico/log-processing-microservices/common/envs"
	"github.com/jose-lico/log-processing-microservices/common/kafka"
	"github.com/jose-lico/log-processing-microservices/common/logging"
	pb "github.com/jose-lico/log-processing-microservices/common/protos"

	log_types "github.com/jose-lico/log-processing-microservices/common/types"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/IBM/sarama"
	"go.uber.org/zap"
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

	kafkaHost := os.Getenv("KAFKA_HOST")
	kafkaPort := os.Getenv("KAFKA_PORT")

	consumerGroup, err := kafka.CreateKafkaConsumer([]string{fmt.Sprintf("%s:%s", kafkaHost, kafkaPort)}, "logs_group")
	if err != nil {
		logging.Logger.Fatal("Failed to start Kafka producer", zap.Error(err))
	}
	logging.Logger.Info("Started kakfa producer")
	defer consumerGroup.Close()

	ctx, cancel := context.WithCancel(context.Background())
	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, os.Interrupt)

	consumer := Consumer{}
	go func() {
		for {
			if err := consumerGroup.Consume(ctx, []string{"logs"}, &consumer); err != nil {
				logging.Logger.Error("Error from consumer", zap.Error(err))
			}
			if ctx.Err() != nil {
				return
			}
		}
	}()

	<-sigterm
	cancel()
	logging.Logger.Info("Shutting down consumer")
}

type Consumer struct{}

func (c *Consumer) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

func (c *Consumer) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (c *Consumer) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		logging.Logger.Info("Log Message received", zap.String("value", string(msg.Value)), zap.String("topic", msg.Topic), zap.Int32("partition", msg.Partition), zap.Int64("offset", msg.Offset))

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

		entry := &pb.ProcessLogEntry{
			Timestamp:      logEntry.Timestamp,
			Level:          logEntry.Level,
			Message:        logEntry.Message,
			UserId:         logEntry.UserID,
			AdditionalData: logEntry.AdditionalData,
			Processed:      logEntry.Processed,
		}

		response, err := client.SubmitLog(ctx, entry)
		if err != nil {
			logging.Logger.Error("Error submitting log", zap.Error(err))
			return err
		}

		logging.Logger.Info("Received storage response", zap.String("response", response.GetMessage()))

		// Consume message regardless for now, handle status after store is done

		sess.MarkMessage(msg, "")

		logging.Logger.Info("Log Message consumed", zap.String("value", string(msg.Value)), zap.String("topic", msg.Topic), zap.Int32("partition", msg.Partition), zap.Int64("offset", msg.Offset))
	}
	return nil
}
