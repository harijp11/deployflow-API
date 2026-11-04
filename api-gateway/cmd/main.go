package main

import (
	"log"
	"net/http"

	"deployflow/api-gateway/internal/config"
	"deployflow/api-gateway/internal/initializer"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	server, err := initializer.Initialize(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize gateway: %v", err)
	}

	log.Printf("API Gateway running on port %s", cfg.Port)

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Gateway server failed: %v", err)
	}
}
