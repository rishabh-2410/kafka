# Kafka Producer/Consumer

A Go-based Kafka producer and consumer application demonstrating how to produce and consume messages using the Confluent Kafka Go client library.

## Overview

This project provides two separate applications:
- **Producer**: Generates and sends mock purchase events to a Kafka topic
- **Consumer**: Subscribes to and consumes messages from a Kafka topic

## Features

- ✅ Kafka producer with SASL/SSL authentication
- ✅ Kafka consumer with automatic offset management
- ✅ Mock purchase data generation
- ✅ Docker Compose setup for local Kafka cluster
- ✅ Environment-based configuration
- ✅ Graceful shutdown handling

## Prerequisites

- Go 1.24.5 or higher
- Docker and Docker Compose (for local Kafka setup)
- Access to a Kafka broker (or use the provided Docker setup)

## Project Structure

```
.
├── cmd/
│   ├── producer/
│   │   └── main.go      # Kafka producer application
│   └── consumer/
│       └── main.go      # Kafka consumer application
├── docker-compose.yml   # Local Kafka and Zookeeper setup
├── go.mod               # Go module dependencies
├── go.sum               # Go module checksums
└── README.md            # This file
```

## Installation

1. Clone the repository:
```bash
git clone <repository-url>
cd kafka
```

2. Install Go dependencies:
```bash
go mod download
```

## Configuration

Create a `.env` file in the root directory with the following variables:

```env
BROKER=localhost:9092
KEY=your-kafka-api-key
SECRET=your-kafka-api-secret
```

**For local development** (using Docker Compose), you can use:
```env
BROKER=<localhost:9092
KEY=
SECRET=
```

Note: When using the local Docker setup, KEY and SECRET can be empty as it uses PLAINTEXT protocol.

## Running Locally with Docker

1. Start Kafka and Zookeeper using Docker Compose:
```bash
docker-compose up -d
```

2. Wait for Kafka to be ready (usually takes 10-20 seconds)

3. Verify the services are running:
```bash
docker-compose ps
```

4. To stop the services:
```bash
docker-compose down
```

## Usage

### Running the Producer

The producer generates mock purchase events and sends them to the `purchases` topic:

```bash
go run cmd/producer/main.go
```

The producer will:
- Generate 10 random purchase events
- Use random users as message keys
- Use random items as message values
- Print delivery confirmations for each message

### Running the Consumer

The consumer subscribes to the `purchases` topic and continuously reads messages:

```bash
go run cmd/consumer/main.go
```

The consumer will:
- Subscribe to the `purchases` topic
- Read messages with a 100ms timeout
- Print consumed messages to the console
- Handle graceful shutdown on SIGINT or SIGTERM (Ctrl+C)

### Running Both Together

1. Start the consumer in one terminal:
```bash
go run cmd/consumer/main.go
```

2. In another terminal, run the producer:
```bash
go run cmd/producer/main.go
```

You should see the consumer receiving the messages produced by the producer.

## Building

Build the producer:
```bash
go build -o bin/producer cmd/producer/main.go
```

Build the consumer:
```bash
go build -o bin/consumer cmd/consumer/main.go
```

Run the binaries:
```bash
./bin/producer
./bin/consumer
```

## Dependencies

- [confluent-kafka-go](https://github.com/confluentinc/confluent-kafka-go) - Confluent's Kafka client library for Go
- [godotenv](https://github.com/joho/godotenv) - Environment variable management

## Kafka Configuration

### Producer Settings
- **Security Protocol**: SASL_SSL
- **SASL Mechanism**: PLAIN
- **Acknowledgment**: all (ensures message durability)

### Consumer Settings
- **Security Protocol**: SASL_SSL
- **SASL Mechanism**: PLAIN
- **Auto Offset Reset**: earliest (reads from the beginning if no offset exists)
- **Group ID**: kafka-consumer-groud-0

## Topic Information

- **Topic Name**: `purchases`
- **Message Key**: User ID (e.g., "eabara", "jsmith")
- **Message Value**: Item purchased (e.g., "book", "alarm clock")

## Troubleshooting

### Connection Issues
- Ensure Kafka broker is running and accessible
- Verify BROKER, KEY, and SECRET in `.env` are correct
- Check firewall settings if connecting to remote broker

### Local Docker Issues
- Ensure Docker is running
- Check if ports 9092 and 2181 are available
- View logs: `docker-compose logs kafka`

### Consumer Not Receiving Messages
- Ensure the topic exists (auto-created by default)
- Check consumer group ID is correct
- Verify producer is sending to the same topic
