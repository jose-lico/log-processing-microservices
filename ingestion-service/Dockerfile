FROM golang:1.22 AS builder

WORKDIR /app

COPY ./ingestion-service/go.mod ./ingestion-service/
COPY ./common/go.mod ./common/go.sum ./common/

RUN cd ./ingestion-service && go mod download

COPY ./ingestion-service/ ./ingestion-service/
COPY ./common/ ./common/

RUN cd ./ingestion-service && CGO_ENABLED=0 GOOS=linux go build -o /app/bin/main ./cmd/main/

FROM alpine:3.19
RUN apk add --no-cache ca-certificates

WORKDIR /app

COPY --from=builder /app/bin/main /app

CMD ["./main"]