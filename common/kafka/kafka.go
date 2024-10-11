package kafka

import (
	"fmt"

	"github.com/jose-lico/log-processing-microservices/common/logging"

	"github.com/IBM/sarama"
	"go.uber.org/zap"
)

func CreateKafkaProducer(brokers []string) (sarama.SyncProducer, error) {
	saramaLogger := &ZapSaramaLogger{
		logger: logging.Logger,
	}
	sarama.Logger = saramaLogger

	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5

	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		return nil, err
	}

	return producer, nil
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
