module github.com/jose-lico/log-processing-microservices/ingestion-service

go 1.22.0

require github.com/jose-lico/log-processing-microservices/common v0.0.0

require (
	github.com/go-chi/chi/v5 v5.1.0 // indirect
	github.com/joho/godotenv v1.5.1 // indirect
	github.com/rs/cors v1.11.1 // indirect
)

replace github.com/jose-lico/log-processing-microservices/common => ../common
