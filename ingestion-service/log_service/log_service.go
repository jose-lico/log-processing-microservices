package log_service

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/jose-lico/log-processing-microservices/common/api"

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

type LogEntry struct {
	Timestamp      string                 `json:"timestamp" validate:"required,datetime=2006-01-02T15:04:05Z07:00"`
	Level          string                 `json:"level" validate:"required,oneof=INFO WARN ERROR DEBUG"`
	Message        string                 `json:"message" validate:"required"`
	UserID         string                 `json:"uesrId" validate:"omitempty,uuid4"`
	AdditionalData map[string]interface{} `json:"additionalData" validate:"omitempty"`
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
		response, _ := json.Marshal(map[string]interface{}{
			"status":  "error",
			"message": "Validation failed.",
			"errors":  errorMessages,
		})
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(response)
		return
	}

	// TODO: Process the log

	err = api.WriteJSON(w, http.StatusOK, map[string]string{
		"status":  "success",
		"message": "Log received successfully.",
	})

	if err != nil {
		fmt.Println(err)
	}
}
