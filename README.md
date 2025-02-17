# Go Lambda Template Service

This is a simple Go service designed to run both as an AWS Lambda function and as a local HTTP server. It uses Gorilla Mux for routing, and AWS Lambda Go packages to bridge the gap between API Gateway events and standard HTTP requests.

## Features

- Health check endpoint (`/` and `/health`) returning the service status and current time.
- Structured logging using Logrus in JSON format.
- Adaptive behavior: runs as an AWS Lambda function or as a local server based on environment variables.

## Prerequisites

- Go 1.16+ installed.
- AWS credentials (if deploying on AWS Lambda).
- Set the environment variable `AWS_LAMBDA_FUNCTION_NAME` for AWS Lambda mode.

## Running Locally

To run the service locally:

```bash
go run ./cmd/service/main.go
```

The service will start on port 8080. You can test the health endpoint by visiting [http://localhost:8080/health](http://localhost:8080/health).

## Running on AWS Lambda

Deploy the service to AWS Lambda. The handler entrypoint is the `Handler` function in `cmd/service/main.go`. AWS Lambda will automatically detect the execution environment.

## Building the Service

To build the service for deployment:

```bash
go build -o service ./cmd/service/main.go
```

## Logging

Logs are output in JSON format. Adjust the log level in the code if necessary.

## License

// ...existing code...
