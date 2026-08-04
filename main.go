package main

import (
	"flag"
	"log"

	"github.com/rahmannugar/proxify/config"
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

	srv := server.New(cfg)

	log.Printf("Starting Proxify on :%d...", cfg.Server.Port)

	if err := srv.Start(); err != nil {
		log.Fatal(err)
	}
}
