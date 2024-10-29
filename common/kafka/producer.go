package kafka

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/jose-lico/log-processing-microservices/common/logging"

	"github.com/IBM/sarama"
	"go.uber.org/zap"
)

const (
	maxRetries        = 5
	reconnectCooldown = 3
)

type Async struct {
	Producer sarama.AsyncProducer
	wg       sync.WaitGroup
}

func NewAsyncProducer(ctx context.Context, id string, brokers []string) (*Async, error) {
	saramaLogger := &ZapSaramaLogger{
		logger: logging.Logger,
	}
	sarama.Logger = saramaLogger

	config := sarama.NewConfig()
	config.ClientID = id
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5

	var err error

	for attempts := 1; attempts <= maxRetries; attempts++ {
		select {
		case <-ctx.Done():
			logging.Logger.Info("Kafka producer creation canceled")

			return nil, ctx.Err()
		default:
			producer, err := sarama.NewAsyncProducer(brokers, config)
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

			kp := &Async{
				Producer: producer,
			}

			kp.wg.Add(2)
			go kp.handleSuccess()
			go kp.handleError()

			logging.Logger.Info("Kafka producer created successfully", zap.String("id", id))

			return kp, nil
		}
	}

	logging.Logger.Error(
		fmt.Sprintf("Failed to create Kafka producer after %d attempts", maxRetries),
		zap.Error(err),
	)

	return nil, err
}

func (kp *Async) Close() {
	kp.Producer.AsyncClose()
	kp.wg.Wait()
	logging.Logger.Info("Kafka producer closed")
}

func (kp *Async) handleSuccess() {
	defer kp.wg.Done()
	for msg := range kp.Producer.Successes() {
		logging.Logger.Info("Message sent successfully",
			zap.String("topic", msg.Topic),
			zap.Int32("partition", msg.Partition),
			zap.Int64("offset", msg.Offset))
	}
}

func (kp *Async) handleError() {
	defer kp.wg.Done()
	for err := range kp.Producer.Errors() {
		logging.Logger.Error("Failed to produce message", zap.Error(err))
	}
}
