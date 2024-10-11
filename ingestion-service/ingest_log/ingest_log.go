package ingest_log

import (
	"encoding/json"
	"net/http"

	"github.com/jose-lico/log-processing-microservices/common/api"
	"github.com/jose-lico/log-processing-microservices/common/logging"

	"github.com/IBM/sarama"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

type Service struct {
	kafkaProducer sarama.SyncProducer
}

func NewService(p sarama.SyncProducer) *Service {
	return &Service{kafkaProducer: p}
}

func (s *Service) RegisterRoutes(r chi.Router) {
	r.Post("/", s.ingestLog)
}

func (s *Service) ingestLog(w http.ResponseWriter, r *http.Request) {
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
	if err := s.publishLog(topic, "This is a fake message"); err != nil {
		api.WriteJSON(w, http.StatusInternalServerError, map[string]interface{}{
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
		logging.Logger.Error("Error writing JSON Response to client", zap.Error(err))
	}
}

func (s *Service) publishLog(topic string, message string) error {
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(message),
	}

	partition, offset, err := s.kafkaProducer.SendMessage(msg)
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
