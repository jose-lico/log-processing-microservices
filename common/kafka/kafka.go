package kafka

import (
	"context"
	"fmt"
	"time"

	"github.com/jose-lico/log-processing-microservices/common/logging"

	"github.com/IBM/sarama"
	"go.uber.org/zap"
)

const (
	maxRetries        = 5
	reconnectCooldown = 3
)

func CreateKafkaProducer(ctx context.Context, brokers []string) (sarama.SyncProducer, error) {
	saramaLogger := &ZapSaramaLogger{
		logger: logging.Logger,
	}
	sarama.Logger = saramaLogger

	config := sarama.NewConfig()
	config.ClientID = "ingestion-service"
	config.Producer.Return.Successes = true
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5

	var producer sarama.SyncProducer
	var err error

	for attempts := 1; attempts <= maxRetries; attempts++ {
		select {
		case <-ctx.Done():
			logging.Logger.Info("Kafka producer creation canceled")
			return nil, ctx.Err()
		default:
			producer, err = sarama.NewSyncProducer(brokers, config)
			if err != nil {
				logging.Logger.Warn(
					fmt.Sprintf("Failed to create Kafka producer (attempt %d/%d)", attempts, maxRetries),
					zap.Error(err),
				)
				if attempts < maxRetries {
					time.Sleep(reconnectCooldown * time.Second)
				}
				continue
			}
			logging.Logger.Info("Kafka producer created successfully")
			return producer, nil
		}
	}

	logging.Logger.Error(
		fmt.Sprintf("Failed to create Kafka producer after %d attempts", maxRetries),
		zap.Error(err),
	)

	return nil, err
}

func CreateKafkaConsumer(addrs []string, groupID string) (sarama.ConsumerGroup, error) {
	config := sarama.NewConfig()
	config.Consumer.Group.Rebalance.Strategy = sarama.NewBalanceStrategyRoundRobin()
	config.Consumer.Offsets.Initial = sarama.OffsetOldest

	consumerGroup, err := sarama.NewConsumerGroup(addrs, groupID, config)
	if err != nil {
		return nil, err
	}

	return consumerGroup, nil
}

type ZapSaramaLogger struct {
	logger *zap.Logger
}

func (l *ZapSaramaLogger) Print(v ...interface{}) {
	l.logger.Info(fmt.Sprint(v...))
}

func (l *ZapSaramaLogger) Printf(format string, v ...interface{}) {
	l.logger.Info(fmt.Sprintf(format, v...))
}

func (l *ZapSaramaLogger) Println(v ...interface{}) {
	l.logger.Info(fmt.Sprintln(v...))
}
