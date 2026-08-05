package main

import (
	"context"
	"flag"
	"log"

	"github.com/rahmannugar/proxify/config"
	"github.com/rahmannugar/proxify/health"
	"github.com/rahmannugar/proxify/loadbalancer"
	"github.com/rahmannugar/proxify/middleware"
	"github.com/rahmannugar/proxify/proxy"
	"github.com/rahmannugar/proxify/registry"
	"github.com/rahmannugar/proxify/server"
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

	//init registry
	registry, err := registry.New(cfg.Upstreams)
	if err != nil {
		log.Fatalf("failed to create registry: %v", err)
	}

	//init health checker
	healthChecker := health.New(
		registry,
		cfg.Health.Endpoint,
		cfg.Health.Interval,
		cfg.Health.Timeout,
		cfg.Health.ExpectedStatus,
	)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go healthChecker.Start(ctx)

	//init load balancer
	loadBalancer := loadbalancer.NewRoundRobin(registry)

	//init proxy
	proxy := proxy.New(loadBalancer)

	//init middleware
	handler := middleware.Logging(proxy)

	//init server
	srv := server.New(cfg, handler)
	log.Printf("Starting Proxify on :%d...", cfg.Server.Port)

	if err := srv.Start(); err != nil {
		log.Fatal(err)
	}
}
