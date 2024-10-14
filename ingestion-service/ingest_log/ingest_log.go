package ingest_log

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/jose-lico/log-processing-microservices/common/api"
	"github.com/jose-lico/log-processing-microservices/common/logging"
	"github.com/jose-lico/log-processing-microservices/common/types"

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
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		api.WriteJSON(w, http.StatusInternalServerError, map[string]interface{}{
			"status":  "error",
			"message": "Unable to read request body",
			"error":   err,
		})
		return
	}
	defer r.Body.Close()

	var logEntry types.IngestLogEntry

	err = json.Unmarshal(bodyBytes, &logEntry)
	if err != nil {
		api.WriteJSON(w, http.StatusInternalServerError, map[string]interface{}{
			"status":  "error",
			"message": "Unable unmarshal data to JSON",
			"error":   err,
		})
		return
	}

	err = api.Validate.Struct(logEntry)
	if err != nil {
		validationErrors := err.(validator.ValidationErrors)
		var errorMessages []string
		for _, ve := range validationErrors {
			errorMessages = append(errorMessages, fmt.Sprintf("%s is invalid", ve.Field()))
		}

		api.WriteJSON(w, http.StatusBadRequest, map[string]interface{}{
			"status":  "error",
			"message": "Validation failed.",
			"errors":  errorMessages,
		})
		return
	}

	topic := "logs"
	if err := s.publishLog(topic, bodyBytes); err != nil {
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

func (s *Service) publishLog(topic string, message []byte) error {
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.ByteEncoder(message),
	}

	partition, offset, err := s.kafkaProducer.SendMessage(msg)
	if err != nil {
		return err
	}

	logging.Logger.Info("Log message published",
		zap.String("topic", topic),
		zap.ByteString("message", message),
		zap.Int32("partition", partition),
		zap.Int64("offset", offset))
	return nil
}
