package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/IBM/sarama"
	"github.com/jose-lico/log-processing-microservices/common/envs"
	"github.com/jose-lico/log-processing-microservices/common/kafka"
	"github.com/jose-lico/log-processing-microservices/common/logging"

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
		log.Printf("Message claimed: value = %s, topic = %s, partition = %d, offset = %d",
			string(msg.Value), msg.Topic, msg.Partition, msg.Offset)

		sess.MarkMessage(msg, "")
	}
	return nil
}
