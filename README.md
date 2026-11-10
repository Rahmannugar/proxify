# Proxify

Proxify is a Layer 7 (application-layer) HTTP load balancer and reverse proxy written in Go that sits between clients and upstream services. It selects an upstream for each incoming HTTP request using round-robin or least-outstanding load balancing, then forwards the request and response.

Proxify operates at the HTTP request level and provides health checks, circuit breaking, rate limiting, request logging, metrics, and graceful shutdown.

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
- YAML

## Running Locally

Start the example upstream services in separate terminals.

Terminal 1:

```bash
go run ./examples/upstream \
  --name api-1 \
  --port 9001
```

Terminal 2:

```bash
go run ./examples/upstream \
  --name api-2 \
  --port 9002
```

Then start Proxify:

```bash
go run ./cmd/proxify
```

or with a specific configuration:

```bash
go run ./cmd/proxify --config configs/least-outstanding.yaml
```

## Running with Docker

Start the entire environment with:

```bash
docker compose up --build
```

Docker Compose starts:

- Proxify
- Upstream API 1
- Upstream API 2

## Using Proxify with Your Own Services

Proxify can also be used as a standalone reverse proxy for your own applications.

Start with one of the example configuration files in the `configs/` directory and update the upstreams to point to your own services.

```yaml
upstreams:
  - id: api-1
    url: http://api-1:9001

  - id: api-2
    url: http://api-2:9002
```

Then start Proxify with your configuration:

```bash
./proxify --config config.yaml
```

or

```bash
go run ./cmd/proxify --config config.yaml
```

If you're using the Docker image:

```bash
docker run \
  -p 8080:8080 \
  -v $(pwd)/config.yaml:/app/config.yaml \
  rahmannugar/proxify \
  --config /app/config.yaml
```

Any compatible YAML configuration can be supplied using the `--config` flag.

## Example Configurations

The repository includes several example configurations:

- `configs/round-robin.yaml`
- `configs/least-outstanding.yaml`
- `configs/docker.yaml`

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
