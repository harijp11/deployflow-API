package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Port        string
	Environment string

	AuthServiceURL       string
	UserServiceURL       string
	ProjectServiceURL    string
	DeploymentServiceURL string
	BuildServiceURL      string
	RuntimeServiceURL    string

	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

func Load() (*Config, error) {
	// Load .env if available.
	// Environment variables already set in the OS take precedence.
	_ = godotenv.Load()

	cfg := &Config{
		Port:                 getEnv("GATEWAY_PORT", "8080"),
		Environment:          getEnv("APP_ENV", "development"),
		AuthServiceURL:       getEnv("AUTH_SERVICE_URL", "http://localhost:8081"),
		UserServiceURL:       getEnv("USER_SERVICE_URL", "http://localhost:8082"),
		ProjectServiceURL:    getEnv("PROJECT_SERVICE_URL", "http://localhost:8083"),
		DeploymentServiceURL: getEnv("DEPLOYMENT_SERVICE_URL", "http://localhost:8084"),
		BuildServiceURL:      getEnv("BUILD_SERVICE_URL", "http://localhost:8085"),
		RuntimeServiceURL:    getEnv("RUNTIME_SERVICE_URL", "http://localhost:8086"),
		ReadTimeout:          15 * time.Second,
		WriteTimeout:         30 * time.Second,
	}

	if cfg.Port == "" {
		return nil, fmt.Errorf("gateway port cannot be empty")
	}

	if _, err := strconv.Atoi(cfg.Port); err != nil {
		return nil, fmt.Errorf("invalid gateway port: %s", cfg.Port)
	}

	return cfg, nil
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	return value
}
