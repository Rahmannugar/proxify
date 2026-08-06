package main

import (
	"context"
	"flag"
	"log"

	"github.com/rahmannugar/proxify/internal/config"
	"github.com/rahmannugar/proxify/internal/exporter"
	"github.com/rahmannugar/proxify/internal/health"
	"github.com/rahmannugar/proxify/internal/loadbalancer"
	"github.com/rahmannugar/proxify/internal/metrics"
	"github.com/rahmannugar/proxify/internal/middleware"
	"github.com/rahmannugar/proxify/internal/proxy"
	"github.com/rahmannugar/proxify/internal/ratelimiter"
	"github.com/rahmannugar/proxify/internal/registry"
	"github.com/rahmannugar/proxify/internal/server"
	"github.com/rahmannugar/proxify/internal/transport"
)

const defaultConfigPath = "config.yaml"

func main() {
	// Allow a custom configuration file to be supplied at startup.
	configPath := flag.String(
		"config",
		defaultConfigPath,
		"Path to the configuration file.",
	)

	flag.Parse()

	cfg, err := config.Load(*configPath)
	if err != nil {
		log.Fatalf("failed to load configuration: %v", err)
	}

	// Init shared transport.
	sharedTransport := transport.New(cfg.Transport)

	// Init registry.
	registry, err := registry.New(
		cfg.Upstreams,
		cfg.CircuitBreaker,
		sharedTransport,
	)
	if err != nil {
		log.Fatalf("failed to create registry: %v", err)
	}

	// Init metrics collector.
	metricsCollector := metrics.New()

	// Init metrics exporter.
	metricsExporter := exporter.NewJSONMetrics(metricsCollector)

	// Init health checker.
	healthChecker := health.New(
		registry,
		metricsCollector,
		cfg.Health.Endpoint,
		cfg.Health.Interval,
		cfg.Health.Timeout,
		cfg.Health.ExpectedStatus,
	)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go healthChecker.Start(ctx)

	// Init load balancer.
	loadBalancer, err := loadbalancer.NewLoadBalancer(
		cfg.LoadBalancer.Strategy,
		registry,
	)
	if err != nil {
		log.Fatalf("failed to create load balancer: %v", err)
	}

	// Init proxy.
	proxy := proxy.New(
		loadBalancer,
		metricsCollector,
	)

	// Init rate limiter.
	limiter := ratelimiter.New(
		cfg.RateLimit.RequestsPerSecond,
		cfg.RateLimit.Burst,
	)

	// Init middleware.
	handler := middleware.RequestID(
		middleware.Logging(
			middleware.RateLimit(
				limiter,
				metricsCollector,
				proxy,
			),
		),
	)

	// Init server.
	srv := server.New(
		cfg,
		handler,
		metricsExporter,
	)

	log.Printf("Starting Proxify on :%d...", cfg.Server.Port)

	if err := srv.Start(); err != nil {
		log.Fatal(err)
	}
}
