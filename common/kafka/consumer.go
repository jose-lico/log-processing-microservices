package kafka

import (
	"github.com/jose-lico/log-processing-microservices/common/logging"

	"github.com/IBM/sarama"
	"go.uber.org/zap"
)

func CreateKafkaConsumer(addrs []string, id, groupID string) (sarama.ConsumerGroup, error) {
	saramaLogger := &ZapSaramaLogger{
		logger: logging.Logger,
	}
	sarama.Logger = saramaLogger

	config := sarama.NewConfig()
	config.ClientID = id
	config.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategyRoundRobin()}
	config.Consumer.Offsets.Initial = sarama.OffsetOldest

	consumerGroup, err := sarama.NewConsumerGroup(addrs, groupID, config)
	if err != nil {
		return nil, err
	}

	logging.Logger.Info("Kafka consumer created successfully", zap.String("id", id))

	return consumerGroup, nil
}
