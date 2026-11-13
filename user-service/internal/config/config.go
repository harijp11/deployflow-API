package config

import (
    "fmt"
    "os"
    "strings"

    "github.com/joho/godotenv"
)

type Config struct {
    App      AppConfig
    Server   ServerConfig
    Database DatabaseConfig
}

type AppConfig struct {
    Env     string
    GinMode string
}

type ServerConfig struct {
    Port string
}

type DatabaseConfig struct {
    URL string
}

func Load() (*Config, error) {
    _ = godotenv.Load("user-service/.env")

    port := os.Getenv("USER_SERVICE_PORT")
    if port == "" {
        port = "8082"
    }

    cfg := &Config{
        App: AppConfig{
            Env:     getEnv("APP_ENV", "development"),
            GinMode: getEnv("GIN_MODE", "debug"),
        },
        Server: ServerConfig{
            Port: port,
        },
        Database: DatabaseConfig{
            URL: os.Getenv("DATABASE_URL"),
        },
    }

    if cfg.Database.URL == "" {
        return nil, fmt.Errorf("DATABASE_URL is required")
    }

    return cfg, nil
}

func getEnv(key, fallback string) string {
    value := strings.TrimSpace(os.Getenv(key))
    if value == "" {
        return fallback
    }
    return value
}