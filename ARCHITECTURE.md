# Architecture

## Overview

Proxify is a reverse proxy that sits between clients and one or more upstream services.

Clients communicate with Proxify instead of communicating directly with upstream services. For every incoming request, Proxify selects a healthy upstream using the configured load-balancing strategy, forwards the request, receives the upstream response, and returns that response to the client.

```text
                Client
                   │
                   ▼
            +--------------+
            |   Proxify    |
            +--------------+
              │    │    │
              ▼    ▼    ▼
         Upstream  Upstream  Upstream
            A         B         C
```

## Responsibilities

Proxify is responsible for:

- Accepting incoming HTTP requests
- Assigning request IDs
- Applying rate limiting
- Logging requests
- Selecting healthy upstreams
- Load balancing
- Circuit breaking
- Forwarding requests to upstream services
- Returning upstream responses
- Performing periodic health checks
- Collecting metrics
- Graceful shutdown

## Non-Responsibilities

Proxify is not responsible for:

- Business logic
- Authentication and authorization
- Data persistence
- Database access
- Service-specific processing

These responsibilities belong to the upstream services.

## Core Components

- HTTP Server
- Reverse Proxy
- Configuration
- Upstream Registry
- Shared HTTP Transport
- Load Balancer
- Circuit Breaker
- Health Checker
- Rate Limiter
- Request ID Middleware
- Logging Middleware
- Metrics Collector

## Request Flow

1. A client sends an HTTP request to Proxify.
2. A request ID is assigned if one is not already present.
3. The request passes through the rate limiter.
4. Request information is logged.
5. Proxify selects a healthy upstream using the configured load-balancing strategy.
6. The circuit breaker determines whether requests may be forwarded to the selected upstream.
7. The request is forwarded to the upstream.
8. The upstream processes the request.
9. The upstream returns a response.
10. Proxify records metrics and forwards the response back to the client.

## Health Checking

A background health checker periodically sends requests to each upstream's configured health endpoint. Healthy upstreams remain eligible for load balancing, while unhealthy upstreams are skipped until they recover.

## Circuit Breaking

Each upstream maintains an independent circuit breaker. When repeated failures exceed the configured threshold, the breaker opens and temporarily blocks requests to that upstream. After the recovery timeout expires, the breaker transitions to a half-open state to determine whether the upstream has recovered.

## Deployment

Proxify can be run either directly on the host or inside Docker.

The Docker Compose environment starts:

- Proxify
- Upstream API 1
- Upstream API 2

allowing the complete request flow, health checks, and load balancing behavior to be tested locally.
