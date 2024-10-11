package ingest_log

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
