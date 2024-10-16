.PHONY: dev run build

# ================================
#          Go Commands
# ================================

dev:
	@echo "Running ingestion-service with live reload..."
	ENV=LOCAL air

run:
	ENV=LOCAL go run ./cmd/main

build:
	go build -o bin/main ./cmd/main
