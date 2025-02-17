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

## TODO

- [x] Create a base code with mux server that can run on Lambda.
- [ ] Create a GitHub CI/CD pipeline for deploying to Lambda.
- [ ] Write a README detailing steps on the AWS side for configuring Lambda.
- [ ] Write a README detailing steps on how to integrate a custom domain with the Lambda function.

## License

MIT License

Copyright (c) [Year] [Your Name]

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
