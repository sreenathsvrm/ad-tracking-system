Ad Tracking System
A scalable and resilient Go backend for managing and tracking video advertisements. This system provides APIs to manage ads, track user clicks, and fetch real-time analytics.

Features
API Endpoints: Manage ads, track clicks, and fetch analytics.

Scalability: Handles high traffic and surges efficiently.

Resilience: No data loss with circuit breakers and retries.

Real-time Analytics: Powered by Redis for fast queries.

Docker Support: Easy containerization.

Kubernetes Deployment: Ready for production.

Table of Contents
Prerequisites

Local Setup

Running with Docker

Kubernetes Deployment

API Endpoints

Testing the APIs

Monitoring and Logging

Prerequisites
Before you begin, ensure you have the following installed:

Go (v1.22 or higher)

Docker and Docker Compose

Kubectl (for Kubernetes deployment)

Postman or any API testing tool

Git

Local Setup
1. Clone the Repository
bash
Copy
git clone https://github.com/sreenathsvrm/ad-tracking-system.git
cd ad-tracking-system
2. Set Up Environment Variables
Create a .env file in the root directory and add the following:

env
Copy
HTTP_PORT=8080
KAFKA_BROKERS=localhost:9092
KAFKA_TOPIC=ad-clicks
REDIS_URL=localhost:6379
DATABASE_URL=postgres://postgres:password@localhost:5432/ad_tracking?sslmode=disable
METRICS_PORT=2112
READ_TIMEOUT=10s
WRITE_TIMEOUT=10s
3. Install Dependencies
bash
Copy
go mod download
4. Start Dependencies
Run the following services using Docker Compose:

bash
Copy
docker-compose up -d
This will start:

PostgreSQL (for ad and click data)

Redis (for real-time analytics)

Kafka (for click event processing)

5. Run the Application
bash
Copy
go run cmd/ad-service/main.go
The application will start on http://localhost:8080.

Running with Docker
1. Build the Docker Image
bash
Copy
docker build -t ad-service .
2. Run the Application with Docker Compose
Update the docker-compose.yml file to include the ad-service:

yaml
Copy
ad-service:
  image: ad-service
  ports:
    - "8080:8080"
  environment:
    HTTP_PORT: ${HTTP_PORT}
    KAFKA_BROKERS: ${KAFKA_BROKERS}
    KAFKA_TOPIC: ${KAFKA_TOPIC}
    REDIS_URL: ${REDIS_URL}
    DATABASE_URL: ${DATABASE_URL}
    METRICS_PORT: ${METRICS_PORT}
    READ_TIMEOUT: ${READ_TIMEOUT}
    WRITE_TIMEOUT: ${WRITE_TIMEOUT}
  depends_on:
    kafka:
      condition: service_healthy
    redis:
      condition: service_healthy
    postgres:
      condition: service_healthy
Start all services:

bash
Copy
docker-compose up -d
Kubernetes Deployment
1. Set Up Kubernetes Cluster
Ensure you have a Kubernetes cluster running (e.g., Minikube, GKE, EKS).

2. Apply Kubernetes Manifests
Deploy the application using the provided Kubernetes manifests:

bash
Copy
kubectl apply -f kubernetes/
This will deploy:

PostgreSQL, Redis, and Kafka as stateful services.

Ad Service as a deployment with a load balancer.

3. Access the Application
Get the external IP of the ad-service:

bash
Copy
kubectl get svc ad-service
Access the API at http://<EXTERNAL_IP>:80.

API Endpoints
1. Fetch All Ads
Method: GET

URL: /ads

Response:

json
Copy
[
  {
    "id": "1",
    "image_url": "http://example.com/image1.jpg",
    "target_url": "http://example.com/target1"
  }
]
2. Record a Click
Method: POST

URL: /ads/click

Request Body:

json
Copy
{
  "ad_id": "1",
  "playback_time": 30
}
Response:

json
Copy
{
  "status": "Click recorded"
}
3. Fetch Analytics
Method: GET

URL: /ads/analytics?ad_id=1

Response:

json
Copy
{
  "ad_id": "1",
  "click_count": 10
}
Testing the APIs
Using Postman
Import the provided Postman collection (available in the repository).

Set the base URL to http://localhost:8080 (or the Kubernetes external IP).

Test the endpoints:

Fetch all ads.

Record a click.

Fetch analytics.

Using cURL
Fetch all ads:

bash
Copy
curl -X GET http://localhost:8080/ads
Record a click:

bash
Copy
curl -X POST http://localhost:8080/ads/click -d '{"ad_id": "1", "playback_time": 30}'
Fetch analytics:

bash
Copy
curl -X GET http://localhost:8080/ads/analytics?ad_id=1
Monitoring and Logging
1. Prometheus Metrics
Access metrics at http://localhost:2112/metrics.

2. Grafana Dashboards
Set up Grafana to visualize Prometheus metrics.

3. Structured Logging
Logs are output in JSON format for easy parsing.

Troubleshooting
Database Connection Issues: Ensure PostgreSQL is running and the connection string is correct.

Kafka Errors: Check if Kafka brokers are reachable.

Redis Errors: Verify Redis is running and the URL is correct.

Contributing
Feel free to open issues or submit pull requests for improvements.

License
This project is licensed under the MIT License. See the LICENSE file for details
