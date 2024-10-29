.PHONY: dev up down build k8s-deploy k8s-delete

INGESTION_DIR=./ingestion-service
KAFKA_DIR=./kafka
PROCESSING_DIR=./processing-service
STORAGE_DIR=./storage-service
QUERY_DIR=./query-service

# ================================
#        Local Development
# ================================

dev:
	@echo "Running all services with live reload..."
	@trap '$(MAKE) -C $(KAFKA_DIR) down; exit' INT; \
	$(MAKE) -C $(KAFKA_DIR) up & \
	$(MAKE) -C $(INGESTION_DIR) dev & \
	$(MAKE) -C $(PROCESSING_DIR) dev & \
	$(MAKE) -C $(STORAGE_DIR) dev & \
	$(MAKE) -C $(QUERY_DIR) dev; \
	wait

test:
	@echo "Stress testing with k6"
	k6 run stress-test.js

# ================================
#              Docker
# ================================

up:
	@echo "Spinning up services with docker-compose"
	docker-compose down
	docker-compose build
	docker-compose up -d

down:
	@echo "Spinning down services"
	docker-compose down

build:
	@echo "Building all Docker images..."
	$(MAKE) -C $(INGESTION_DIR) docker-build
#	$(MAKE) -C $(PROCESSING_DIR) docker-build
#	$(MAKE) -C $(STORAGE_DIR) docker-build
#	$(MAKE) -C $(QUERY_DIR) docker-build

# ================================
#            Kubernetes  
# ================================

k8s-deploy-all:
	@echo "Deploying all Kubernetes resources..."
	$(MAKE) -C $(INGESTION_DIR) k8s-deploy
	$(MAKE) -C $(KAFKA_DIR) k8s-deploy
#	$(MAKE) -C $(PROCESSING_DIR) k8s-deploy
#	$(MAKE) -C $(STORAGE_DIR) k8s-deploy
#	$(MAKE) -C $(QUERY_DIR) k8s-deploy

k8s-delete-all:
	@echo "Deleting all Kubernetes resources..."
	$(MAKE) -C $(INGESTION_DIR) k8s-delete
	$(MAKE) -C $(KAFKA_DIR) k8s-delete
#	$(MAKE) -C $(PROCESSING_DIR) k8s-delete
#	$(MAKE) -C $(STORAGE_DIR) k8s-delete
#	$(MAKE) -C $(QUERY_DIR) k8s-delete
