.PHONY: dev docker-build-all k8s-deploy-all k8s-delete-all build-deploy-all

INGESTION_DIR=./ingestion-service
PROCESSING_DIR=./processing-service
STORAGE_DIR=./storage-service
QUERY_DIR=./query-service

# ================================
#          Go Commands
# ================================

dev:
	$(MAKE) -C $(INGESTION_DIR) dev
#	$(MAKE) -C $(PROCESSING_DIR) docker-build
#	$(MAKE) -C $(STORAGE_DIR) docker-build
#	$(MAKE) -C $(QUERY_DIR) docker-build

# ================================
#            Kubernetes  
# ================================

docker-build-all:
	@echo "Building all Docker images..."
	$(MAKE) -C $(INGESTION_DIR) docker-build
#	$(MAKE) -C $(PROCESSING_DIR) docker-build
#	$(MAKE) -C $(STORAGE_DIR) docker-build
#	$(MAKE) -C $(QUERY_DIR) docker-build

# Deploy all Kubernetes resources
k8s-deploy-all:
	@echo "Deploying all Kubernetes resources..."
	$(MAKE) -C $(INGESTION_DIR) k8s-deploy
#	$(MAKE) -C $(PROCESSING_DIR) k8s-deploy
#	$(MAKE) -C $(STORAGE_DIR) k8s-deploy
#	$(MAKE) -C $(QUERY_DIR) k8s-deploy

k8s-delete-all:
	@echo "Deleting all Kubernetes resources..."
	$(MAKE) -C $(INGESTION_DIR) k8s-delete
#	$(MAKE) -C $(PROCESSING_DIR) k8s-delete
#	$(MAKE) -C $(STORAGE_DIR) k8s-delete
#	$(MAKE) -C $(QUERY_DIR) k8s-delete
