.PHONY: k8s-deploy k8s-delete
# ================================
#            Kubernetes  
# ================================

KAFKA_SERVICE_FILE=kafka-service.yml
KAFKA_DEPLOYMENT_FILE=kafka-deployment.yml
ZOOKEEPER_SERVICE_FILE=zookeeper-service.yml
ZOOKEEPER_DEPLOYMENT_FILE=zookeeper-deployment.yml

k8s-deploy:
	@echo "Deploying Kafka to Kubernetes..."
	kubectl apply -f $(KAFKA_SERVICE_FILE)
	kubectl apply -f $(KAFKA_DEPLOYMENT_FILE)
	kubectl apply -f $(ZOOKEEPER_SERVICE_FILE)
	kubectl apply -f $(ZOOKEEPER_DEPLOYMENT_FILE)

k8s-delete:
	@echo "Deleting Kafka Kubernetes resources..."
	kubectl delete -f $(KAFKA_SERVICE_FILE)
	kubectl delete -f $(KAFKA_DEPLOYMENT_FILE)
	kubectl delete -f $(ZOOKEEPER_SERVICE_FILE)
	kubectl delete -f $(ZOOKEEPER_DEPLOYMENT_FILE)