# Video Encoder Service

This service is designed to encode video files using **Golang**, **Minio**, and **Kafka**. It leverages a microservices architecture to provide a scalable and efficient way to process and encode video files stored in Minio, while Kafka is used for communication between the components.

## Table of Contents

- [Overview](#overview)
- [Architecture](#architecture)
- [Getting Started](#getting-started)
  - [Prerequisites](#prerequisites)
  - [Setup](#setup)
  - [Running the Service](#running-the-service)
- [How It Works](#how-it-works)
- [Service Components](#service-components)
  - [Video Encoder](#video-encoder)
  - [Kafka Producer](#kafka-producer)
  - [Minio](#minio)
- [API Endpoints](#api-endpoints)
- [Testing](#testing)
- [Deployment](#deployment)
- [Contributing](#contributing)
- [License](#license)

## Overview

This service is responsible for encoding videos stored in Minio. It listens to messages from a Kafka queue, retrieves the videos from Minio, encodes them, and stores the encoded videos back to Minio.

### Technologies:
- **Golang** for building the core service
- **Minio** for object storage (S3-compatible)
- **Apache Kafka** for message queue handling
- **PostgreSQL** for metadata storage and job tracking

## Architecture

The system consists of the following components:
1. **Kafka Producer**: Sends video processing jobs to Kafka.
2. **Kafka Consumer**: The service that listens to Kafka topics and processes video encoding tasks.
3. **Minio**: An S3-compatible object storage service that stores the raw and encoded videos.
4. **PostgreSQL**: Stores metadata related to encoding jobs, such as status and timestamps.

## Getting Started

### Prerequisites

Before you can run the service, make sure you have the following installed on your machine:
- **Golang** (v1.16 or higher)
- **Docker** (for running Minio, Kafka, and PostgreSQL containers)
- **Kafka** (can be run via Docker)
- **Minio** (can be run via Docker)
- **PostgreSQL** (for storing metadata)

You will also need the following environment variables set:
- `KAFKA_BROKER`: Kafka broker address (e.g., `localhost:9092`)
- `MINIO_ACCESS_KEY`: Minio access key
- `MINIO_SECRET_KEY`: Minio secret key
- `MINIO_BUCKET`: Minio bucket where videos will be stored
- `POSTGRES_DSN`: PostgreSQL connection string

### Setup

1. **Clone the repository:**

   ```bash
   git clone https://github.com/yourusername/video-encoder-service.git
   cd video-encoder-service
   ```

2. **Start dependencies using Docker Compose:**

   ```bash
   docker-compose up -d
   ```

3. **Run the service:**

   ```bash
   go run main.go
   ```

## How It Works

1. A client uploads a raw video to Minio.
2. The client sends a job request to the service via the API.
3. The service publishes a job message to Kafka.
4. The Kafka consumer retrieves the message and processes the video.
5. The encoded video is stored back in Minio.
6. Metadata is saved in PostgreSQL for tracking job status.

## Service Components

### Video Encoder
The core component that processes and encodes videos into different formats.

### Kafka Producer
Responsible for publishing job requests to Kafka for asynchronous processing.

### Minio
An S3-compatible storage solution used for storing raw and encoded videos.

## API Endpoints

### Submit a Video Encoding Job
```http
POST /api/encode
```
**Request Body:**
```json
{
  "video_url": "s3://bucket-name/video.mp4",
  "output_formats": ["mp4", "webm"]
}
```

### Get Encoding Job Status
```http
GET /api/jobs/{job_id}
```

## Testing

Run unit tests using:
```bash
go test ./...
```

## Deployment

1. Build the Docker image:
   ```bash
   docker build -t video-encoder-service .
   ```
2. Push to a container registry:
   ```bash
   docker push your-docker-repo/video-encoder-service
   ```
3. Deploy using Kubernetes or Docker Compose.

## Contributing

1. Fork the repository.
2. Create a new branch (`feature/my-feature`).
3. Commit your changes.
4. Push to your branch and create a Pull Request.

## License

This project is licensed under the MIT License.

