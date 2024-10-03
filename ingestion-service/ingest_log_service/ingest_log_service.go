package ingest_log_service

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/jose-lico/log-processing-microservices/common/api"
	"github.com/jose-lico/log-processing-microservices/common/logging"
	"go.uber.org/zap"

	"github.com/IBM/sarama"
	"github.com/go-playground/validator/v10"
)

// Example log

// {
// 	"timestamp": "2024-10-03T15:20:30Z",
// 	"level": "INFO",
// 	"message": "User logged in",
// 	"userId": "12345",
// 	"additionalData": {
// 	  "ipAddress": "192.168.1.1",
// 	  "sessionId": "abcde12345"
// 	}
// }

var kafkaProducer sarama.SyncProducer

type LogEntry struct {
	Timestamp      string                 `json:"timestamp" validate:"required,datetime=2006-01-02T15:04:05Z07:00"`
	Level          string                 `json:"level" validate:"required,oneof=INFO WARN ERROR DEBUG"`
	Message        string                 `json:"message" validate:"required"`
	UserID         string                 `json:"uesrId" validate:"omitempty,uuid4"`
	AdditionalData map[string]interface{} `json:"additionalData" validate:"omitempty"`
}

func init() {
	sarama.Logger = log.New(os.Stdout, "[Sarama] ", log.LstdFlags)
	brokers := []string{"localhost:9092"}
	var err error
	kafkaProducer, err = createKafkaProducer(brokers)
	if err != nil {
		panic(err)
	}
}

func createKafkaProducer(brokers []string) (sarama.SyncProducer, error) {
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

func publishLog(producer sarama.SyncProducer, topic string, message string) error {
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(message),
	}

	partition, offset, err := producer.SendMessage(msg)
	if err != nil {
		return err
	}

	logging.Logger.Info("Log message published",
		zap.String("topic", topic),
		zap.String("message", message),
		zap.Int32("partition", partition),
		zap.Int64("offset", offset))
	return nil
}

func IngestLog(w http.ResponseWriter, r *http.Request) {
	var logEntry LogEntry

	err := json.NewDecoder(r.Body).Decode(&logEntry)
	if err != nil {
		http.Error(w, "Invalid JSON format.", http.StatusBadRequest)
		return
	}

	err = api.Validate.Struct(logEntry)
	if err != nil {
		validationErrors := err.(validator.ValidationErrors)
		var errorMessages []string
		for _, ve := range validationErrors {
			errorMessages = append(errorMessages, ve.Field()+" is invalid.")
		}

		api.WriteJSON(w, http.StatusBadRequest, map[string]interface{}{
			"status":  "error",
			"message": "Validation failed.",
			"errors":  errorMessages,
		})
		return
	}

	topic := "logs"
	if err := publishLog(kafkaProducer, topic, "This is a fake message"); err != nil {
		api.WriteJSON(w, http.StatusBadRequest, map[string]interface{}{
			"status":  "error",
			"message": "Failed to publish log.",
			"error":   err,
		})
		return
	}

	err = api.WriteJSON(w, http.StatusOK, map[string]string{
		"status":  "success",
		"message": "Log received successfully.",
	})

	if err != nil {
		logging.Logger.Error("Error writing JSON to client", zap.Error(err))
	}
}
