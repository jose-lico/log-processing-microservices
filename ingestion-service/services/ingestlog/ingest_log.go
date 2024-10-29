package ingestlog

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/jose-lico/log-processing-microservices/common/api"
	"github.com/jose-lico/log-processing-microservices/common/kafka"
	"github.com/jose-lico/log-processing-microservices/common/logging"
	"github.com/jose-lico/log-processing-microservices/common/types"

	"github.com/IBM/sarama"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

type Service struct {
	kafkaProducer *kafka.AsyncProducer
}

func NewService(p *kafka.AsyncProducer) *Service {
	return &Service{kafkaProducer: p}
}

func (s *Service) RegisterRoutes(r chi.Router) {
	r.Post("/async", s.ingestLog)
}

func (s *Service) ingestLog(w http.ResponseWriter, r *http.Request) {
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		api.WriteJSON(w, http.StatusInternalServerError, map[string]interface{}{
			"status":  "error",
			"message": "Unable to read request body",
			"error":   err.Error(),
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
			"error":   err.Error(),
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

	msg := &sarama.ProducerMessage{
		Topic: "logs",
		Value: sarama.ByteEncoder(bodyBytes),
	}

	s.kafkaProducer.Producer.Input() <- msg

	err = api.WriteJSON(w, http.StatusAccepted, map[string]string{
		"status":  "success",
		"message": "Log received successfully.",
	})

	if err != nil {
		logging.Logger.Error("Error writing JSON Response to client", zap.Error(err))
	}
}
