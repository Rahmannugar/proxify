# Proxify

Proxify is a reverse proxy written in Go that sits between clients and upstream services. It forwards incoming HTTP requests while providing load balancing, health checks, circuit breaking, rate limiting, request logging, metrics, and graceful shutdown.

## Features

- Reverse proxy
- Request forwarding
- Response forwarding
- Configurable upstreams
- Round-robin load balancing
- Least-outstanding load balancing
- Health checks
- Circuit breaker
- Shared HTTP transport
- Rate limiting
- Request ID middleware
- Request logging
- Metrics endpoint
- Graceful shutdown
- Docker support
- YAML-based configuration

## Technology

- Go
- Docker
- Docker Compose
- YAML

## Running Locally

```bash
go run ./cmd/proxify
```

or with a specific configuration:

```bash
go run ./cmd/proxify --config configs/least-outstanding.yaml
```

## Running with Docker

```bash
docker compose up --build
```

The Docker environment starts:

- Proxify
- Upstream API 1
- Upstream API 2

## Configuration

Proxify is configured using YAML files located in the `configs/` directory.

Available configurations:

- `round-robin.yaml`
- `least-outstanding.yaml`
- `docker.yaml`

## Metrics

Metrics are available at:

```text
GET /metrics
```

## Health Checks

Each upstream exposes:

```text
GET /health
```

The proxy periodically checks upstream health and only routes requests to healthy upstreams.
