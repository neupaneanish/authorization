# Authorization Microservice

###### Design and Developed by [Anish Neupane](https://neupaneanish.com.np)

___

## Overview

Distributed Authorization Microservice with Go, gRPC, Envoy, and Valkey for Envoy authz.

---

## Technologies Stack

![Go](https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white)
![gRPC](https://img.shields.io/badge/gRPC-2596BE?style=for-the-badge&logo=trpc&logoColor=white)
![Valkey](https://img.shields.io/badge/Valkey-FF4438?style=for-the-badge&logo=redis&logoColor=white)
![Envoy](https://img.shields.io/badge/Envoy-AC6199?style=for-the-badge&logo=envoyproxy&logoColor=white)
![Docker](https://img.shields.io/badge/Docker-2496ED?style=for-the-badge&logo=docker&logoColor=white)
![GitHubActions](https://img.shields.io/badge/Actions-2088FF?style=for-the-badge&logo=githubactions&logoColor=white)
![Opentelemetry](https://img.shields.io/badge/Opentelemetry-000000?style=for-the-badge&logo=opentelemetry&logoColor=white)
---

## Environments

|     Name      |               Default               |            Options            |
|:-------------:|:-----------------------------------:|:-----------------------------:|
|  VALKEY_URL   |                                     |                               |
| SERVICE_NAME  | `neupaneanish.com.np/authorization` |                               |
|  ENVIRONMENT  |            `development`            | `development` or `production` |
| TELEMETRY_URL |                                     |        gRPC port only         |
|     PORT      |               `50051`               |        `80` to `65535`        |

```dotenv
VALKEY_URL=127.0.0.1:6379
PORT=50051
SERVICE_NAME=neupaneanish.com.np/authorization
ENVIRONMENT=development
TELEMETRY_URL=127.0.0.1:4317
```

___

## Setup, Execution & Testing

```bash
# 1. Clone the core framework engine
git clone https://github.com/neupaneanish/authorization.git
cd authorization

# 2. Execute the tests
go test -v -tags=unit ./...
go test -v -tags=integration ./...

# 3. Launch the local microservice API server
# (Note: Requires an active OpenTelemetry collector instance, e.g., SigNoz)
go run cmd/server/main.go
```

---

## Coverage ~91.10%

> Note: Metrics reflect core application logic after filtering out `main.go`, and test helper suites.

> Coverage is done through real infrastructure Valkey, OpenTelemetry i.e. testcontainers. It doesn't have
> any mocks.

```bash
# Generate coverage
go test -race -tags=unit,integration -coverprofile=coverage.out -coverpkg=./... ./...

# Filter out external boundaries, generated code, and tooling 
grep -v -E "cmd/|/tests/" coverage.out > coverage.clean.out

# Export to interactive HTML for local branch analysis
go tool cover -html=coverage.clean.out -o coverage.clean.html 

# Output statement breakdown to CLI
go tool cover -func=coverage.clean.out 
```

___

## [License](LICENSE)