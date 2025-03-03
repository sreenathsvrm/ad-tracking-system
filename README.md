# Ad Tracking System

A scalable and resilient Go backend for managing and tracking video advertisements. Provides APIs to manage ads, track clicks, and fetch real-time analytics.

## Features

* **API Endpoints:** Manage ads, track clicks, and fetch analytics.
* **Scalability:** Handles high traffic efficiently.
* **Resilience:** No data loss with circuit breakers and retries.
* **Real-time Analytics:** Powered by Redis.
* **Docker Support:** Included in the project.
* **Kubernetes:** Basic configuration provided (not production-ready).

## Prerequisites

* Go (v1.22+)
* Docker and Docker Compose
* Postman or any API testing tool
* Git

## Quick Start

1.  **Clone the Repository**

    ```bash
    git clone [https://github.com/sreenathsvrm/ad-tracking-system.git](https://github.com/sreenathsvrm/ad-tracking-system.git)
    cd ad-tracking-system
    ```

2.  **Set Up Environment**

    * Create a `.env` file:

        ```ini
        HTTP_PORT=8080
        KAFKA_BROKERS=localhost:9092
        KAFKA_TOPIC=ad-clicks
        REDIS_URL=localhost:6379
        DATABASE_URL=postgres://postgres:password@localhost:5432/ad_tracking?sslmode=disable
        METRICS_PORT=2112
        READ_TIMEOUT=10s
        WRITE_TIMEOUT=10s
        ```

    * Start the services using Docker Compose:

        ```bash
        docker-compose up -d
        ```

    * The API will be available at `http://localhost:8080`.

## API Endpoints

1.  **Fetch All Ads**

    * `GET /ads`
    * **Description:** Returns a list of ads with basic metadata (e.g., ID, image URL, target URL).
    * **Response:**

        ```json
        [
          {
            "id": "1",
            "image_url": "[http://example.com/image1.jpg](http://example.com/image1.jpg)",
            "target_url": "[http://example.com/target1](http://example.com/target1)"
          }
        ]
        ```

2.  **Record a Click**

    * `POST /ads/click`
    * **Request Body:**

        ```json
        {
          "ad_id": "1",
          "playback_time": 30
        }
        ```

    * **Response:**

        ```json
        {
          "status": "Click recorded"
        }
        ```

3.  **Fetch Analytics**

    * `GET /ads/analytics?ad_id=1`
    * **Response:**

        ```json
        {
          "ad_id": "1",
          "click_count": 10
        }
        ```

## Testing the APIs Using cURL

* **Fetch all ads:**

    ```bash
    curl -X GET http://localhost:8080/ads
    ```

* **Record a click:**

    ```bash
    curl -X POST http://localhost:8080/ads/click -d '{"ad_id": "1", "playback_time": 30}'
    ```

* **Fetch analytics:**

    ```bash
    curl -X GET http://localhost:8080/ads/analytics?ad_id=1
    ```

## Monitoring

* **Prometheus Metrics:** `http://localhost:2112/metrics`
* **Grafana:** Set up dashboards to visualize metrics.

## Troubleshooting

* **Database Issues:**
    * **Ensure that the database has dummy data:**
        * Assuming your database has the following `ads` table:

            ```sql
            CREATE TABLE ads (
                id VARCHAR(36) PRIMARY KEY,
                image_url TEXT NOT NULL,
                target_url TEXT NOT NULL
            );
            ```

        * Insert dummy data using SQL:

            ```sql
            INSERT INTO ads (id, image_url, target_url) VALUES
            ('1', '[http://example.com/image1.jpg](http://example.com/image1.jpg)', '[http://example.com/target1](http://example.com/target1)'),
            ('2', '[http://example.com/image2.jpg](https://www.google.com/search?q=http://example.com/image2.jpg)', '[http://example.com/target2](https://www.google.com/search?q=http://example.com/target2)'),
            ('3', '[http://example.com/image3.jpg](https://www.google.com/search?q=http://example.com/image3.jpg)', '[http://example.com/target3](https://www.google.com/search?q=http://example.com/target3)');
            ```

        * Ensure PostgreSQL is running and the connection string is correct.
* **Kafka Errors:** Check if Kafka brokers are reachable.
* **Redis Errors:** Verify Redis is running.

## Contributing

* Open issues or submit pull requests for improvements.

License
This project is licensed under the MIT License. See the LICENSE file for details
