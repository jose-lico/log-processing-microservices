package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/jose-lico/log-processing-microservices/common/envs"
	"github.com/jose-lico/log-processing-microservices/common/kafka"
	"github.com/jose-lico/log-processing-microservices/common/logging"
	log_types "github.com/jose-lico/log-processing-microservices/common/types"

	"github.com/IBM/sarama"
	"go.uber.org/zap"
)

func main() {
	env := os.Getenv("ENV")

	if env == "LOCAL" {
		err := envs.LoadEnvs()
		if err != nil {
			panic(err)
		}
	}

	logging.CreateLogger()
	defer logging.Logger.Sync()

	kafkaHost := os.Getenv("KAFKA_HOST")
	kafkaPort := os.Getenv("KAFKA_PORT")

	consumerGroup, err := kafka.CreateKafkaConsumer([]string{fmt.Sprintf("%s:%s", kafkaHost, kafkaPort)}, "logs_group")
	if err != nil {
		logging.Logger.Fatal("Failed to start Kafka producer", zap.Error(err))
	}
	defer consumerGroup.Close()

	ctx, cancel := context.WithCancel(context.Background())
	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, os.Interrupt)

	consumer := Consumer{}
	go func() {
		for {
			if err := consumerGroup.Consume(ctx, []string{"logs"}, &consumer); err != nil {
				log.Printf("Error from consumer: %v", err)
			}
			if ctx.Err() != nil {
				return
			}
		}
	}()

	<-sigterm
	cancel()
	log.Println("Shutting down consumer...")
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
		logging.Logger.Info("Log Message claimed", zap.String("value", string(msg.Value)), zap.String("topic", msg.Topic), zap.Int32("partition", msg.Partition), zap.Int64("offset", msg.Offset))

		var logEntry log_types.ProccessLogEntry

		err := json.Unmarshal([]byte(msg.Value), &logEntry)
		if err != nil {
			logging.Logger.Error("Invalid JSON format", zap.Error(err))
			return err
		}

		// Dummy business logic...
		logEntry.Processed = true

		sess.MarkMessage(msg, "")
	}
	return nil
}
