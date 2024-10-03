.PHONY: k8s-deploy k8s-delete
# ================================
#            Kubernetes  
# ================================

K8S_SERVICE_FILE=kafka-service.yaml

k8s-deploy:
	@echo "Deploying Kafka to Kubernetes..."
	kubectl apply -f $(K8S_SERVICE_FILE)

k8s-delete:
	@echo "Deleting Kafka Kubernetes resources..."
	kubectl delete -f $(K8S_SERVICE_FILE)