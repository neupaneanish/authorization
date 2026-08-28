# Authorization Microservice

###### Design and Developed by [Anish Neupane](https://neupaneanish.com.np)

___

## Overview

Distributed Authorization Microservice with Go, gRPC, Envoy, and Valkey for Envoy authz.

---

## Technologies Stack

| Technology                                                |                                                                                                  | Description                                                       |
|:----------------------------------------------------------|:------------------------------------------------------------------------------------------------:|:------------------------------------------------------------------|
| [**Go**](https://go.dev)                                  |             <img src="https://thesvg.org/icons/go/default.svg" height="12" alt="Go">             | Core application logic                                            |
| [**gRPC**](https://grpc.io)                               |           <img src="https://thesvg.org/icons/grpc/default.svg" height="24" alt="gRPC">           | High-performance RPC framework                                    |
| [**Valkey**](https://valkey.io)                           |         <img src="https://thesvg.org/icons/valkey/default.svg" height="24" alt="Valkey">         | High-performance data structure store                             |
| [**Envoy**](https://www.envoyproxy.io)                    |         <img src="https://thesvg.org/icons/envoy/default.svg" height="24" alt="Valkey">          | API Gateway and Edge Proxy                                        |
| [**Docker**](https://docker.com)                          |         <img src="https://thesvg.org/icons/docker/default.svg" height="24" alt="Docker">         | Containerization and deployment                                   |
| [**Test Containers**](https://testcontainers.com)         | <img src="https://thesvg.org/icons/development-containers/default.svg" height="24" alt="Docker"> | Orchestrates real Valkey Docker instances inside automated tests. |
| [**GitHub Actions**](https://github.com/features/actions) | <img src="https://thesvg.org/icons/github-actions/default.svg" height="24" alt="GitHub Actions"> | CI/CD automation pipelines                                        |
| [**OpenTelemetry**](https://opentelemetry.io)             |  <img src="https://thesvg.org/icons/opentelemetry/default.svg" height="24" alt="OpenTelemetry">  | Observability and telemetry framework                             |

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