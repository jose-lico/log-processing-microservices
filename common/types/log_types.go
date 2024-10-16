package types

// Example log

// {
// 	"timestamp": "2024-10-03T15:20:30Z",
// 	"level": "INFO",
// 	"message": "User logged in",
// 	"userId": "97a569df-498c-49a3-9676-ec82b992dc07",
// 	"additionalData": {
// 	  "ipAddress": "192.168.1.1",
// 	  "sessionId": "abcde12345"
// 	}
// }

type IngestLogEntry struct {
	Timestamp      string            `json:"timestamp" validate:"required,datetime=2006-01-02T15:04:05Z07:00"`
	Level          string            `json:"level" validate:"required,oneof=INFO WARN ERROR DEBUG"`
	Message        string            `json:"message" validate:"required"`
	UserID         string            `json:"userId" validate:"omitempty,uuid4"`
	AdditionalData map[string]string `json:"additionalData" validate:"omitempty"`
}

type ProccessLogEntry struct {
	Timestamp      string            `json:"timestamp" validate:"required,datetime=2006-01-02T15:04:05Z07:00"`
	Level          string            `json:"level" validate:"required,oneof=INFO WARN ERROR DEBUG"`
	Message        string            `json:"message" validate:"required"`
	UserID         string            `json:"userId" validate:"omitempty,uuid4"`
	AdditionalData map[string]string `json:"additionalData" validate:"omitempty"`

	Processed bool `json:"processed"`
}
