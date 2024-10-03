.PHONY: dev run build docker-build k8s-deploy k8s-delete deploy
# ================================
#          Go Commands
# ================================

dev:
	@echo "Running with live reload..."
	ENV=LOCAL air

run:
	ENV=LOCAL go run ./cmd/main

build:
	go build -o bin/main ./cmd/main

# ================================
#            Kubernetes  
# ================================

IMAGE_NAME=ingestion-service
IMAGE_TAG=local
K8S_DIR=./k8s
K8S_DEPLOYMENT_FILE=$(K8S_DIR)/ingestion-deployment.yaml
K8S_SERVICE_FILE=$(K8S_DIR)/ingestion-service.yaml

docker-build:
	@echo "Building Docker image..."
	cd .. && docker build -t $(IMAGE_NAME):$(IMAGE_TAG) -f ./ingestion-service/Dockerfile .

k8s-deploy: docker-build
	@echo "Deploying to Kubernetes..."
	kubectl apply -f $(K8S_DEPLOYMENT_FILE)
	kubectl apply -f $(K8S_SERVICE_FILE)
	kubectl rollout restart deployment ingestion-deployment

k8s-delete:
	@echo "Deleting Kubernetes resources..."
	kubectl delete -f $(K8S_DEPLOYMENT_FILE)
	kubectl delete -f $(K8S_SERVICE_FILE)
