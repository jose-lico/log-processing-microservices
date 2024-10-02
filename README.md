This project is my attempt at developing a more advanced Go backend application, designed to demonstrate a microservices architecture for log processing.

The system consists of four containerized microservices orchestrated using Kubernetes, utilizing technologies like gRPC, Kafka, and RESTful APIs.
The project aims to simulate a real-world distributed system, providing hands-on experience with interservice communication, container orchestration, authentication, and monitoring.

## Architecture

The system is composed of the following services:

- **Ingestion Service:** Receives logs via RESTful API from clients.
- **Processing Service:** Consumes logs from Kafka, processes them (placeholder logic for now), and communicates with the Storage Service via gRPC.
- **Storage Service:** Handles data persistence using PostgreSQL or Redis (still TBD).
- **Query Service:** Provides a RESTful API for clients to retrieve logs, communicating with the Storage Service via gRPC.
