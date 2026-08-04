# Architecture

## Overview

Proxify is a reverse proxy that sits between clients and one or more upstream services.

Clients communicate with Proxify instead of communicating directly with upstream services. For every incoming request, Proxify selects an upstream based on the configured load-balancing strategy, forwards the request, receives the upstream response, and returns that response to the client.

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
- Selecting an upstream
- Forwarding requests to upstream services
- Returning upstream responses to clients
- Load balancing
- Health checks
- Rate limiting
- Request logging
- Metrics collection
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
- Load Balancer
- Health Checker
- Rate Limiter
- Logger
- Metrics

## Request Flow

1. A client sends an HTTP request to Proxify.
2. Proxify selects a healthy upstream.
3. The request is forwarded to the selected upstream.
4. The upstream processes the request.
5. The upstream returns a response.
6. Proxify forwards the response back to the client.
