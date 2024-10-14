module github.com/jose-lico/log-processing-microservices/storage-service

go 1.22.0

require google.golang.org/grpc v1.67.1

require (
	github.com/jose-lico/log-processing-microservices/common v0.0.0
	go.uber.org/zap v1.27.0
)

require (
	github.com/joho/godotenv v1.5.1 // indirect
	go.uber.org/multierr v1.10.0 // indirect
	golang.org/x/net v0.28.0 // indirect
	golang.org/x/sys v0.24.0 // indirect
	golang.org/x/text v0.17.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240814211410-ddb44dafa142 // indirect
	google.golang.org/protobuf v1.35.1 // indirect
)

replace github.com/jose-lico/log-processing-microservices/common => ../common
