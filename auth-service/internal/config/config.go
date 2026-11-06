
package config

import (
	"fmt"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	App      AppConfig
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
}

type AppConfig struct {
    Name        string
    Environment string
    GinMode     string
}

type ServerConfig struct {
	Port            string
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	IdleTimeout     time.Duration
	ShutdownTimeout time.Duration
	ReadHeaderTimeout time.Duration
}

type DatabaseConfig struct {
	URL             string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
	ConnMaxIdleTime time.Duration
}

type AuthServiceConfig struct {
	Issuer                  string
	Audience                string
	SigningKey              string
	AccessTokenTTL          time.Duration
	RefreshTokenTTL         time.Duration
	PasswordResetTokenTTL   time.Duration
	EmailVerificationTTL    time.Duration
}

func Load() (*Config, error) {
	// Load the local .env only for local development.
	// In production, inject environment variables through the deployment platform.
	if os.Getenv("APP_ENV") != "production" {
		_ = godotenv.Load("auth-service/.env")
	}

	cfg := &Config{
		App: AppConfig{
			Name:        "auth-service",
			Environment: getEnv("APP_ENV", "development"),
		},
		Server: ServerConfig{
			Port:              getEnv("AUTH_SERVICE_PORT", "8081"),
			ReadTimeout:       10 * time.Second,
			ReadHeaderTimeout: 5 * time.Second,
			WriteTimeout:      15 * time.Second,
			IdleTimeout:       60 * time.Second,
			ShutdownTimeout:   10 * time.Second,
		},
		Database: DatabaseConfig{
			URL:             os.Getenv("DATABASE_URL"),
			MaxOpenConns:    25,
			MaxIdleConns:    5,
			ConnMaxLifetime: 30 * time.Minute,
			ConnMaxIdleTime: 5 * time.Minute,
		},
	}

	if err := cfg.validate(); err != nil {
		return nil, err
	}

	jwtConfig, err := loadJWTConfig()
	if err != nil {
		return nil, err
	}
	cfg.JWT = *jwtConfig

	return cfg, nil
}

func (c *Config) validate() error {
	port, err := strconv.Atoi(c.Server.Port)
	if err != nil || port < 1 || port > 65535 {
		return fmt.Errorf("invalid AUTH_SERVICE_PORT")
	}

	switch c.App.Environment {
	case "development", "test", "production":
	default:
		return fmt.Errorf("APP_ENV must be development, test, or production")
	}

	if c.Database.URL == "" {
		return fmt.Errorf("DATABASE_URL is required")
	}

	dbURL, err := url.Parse(c.Database.URL)
	if err != nil {
		return fmt.Errorf("invalid DATABASE_URL: %w", err)
	}

	if (dbURL.Scheme != "postgres" && dbURL.Scheme != "postgresql") ||
		dbURL.Host == "" ||
		dbURL.Path == "" ||
		dbURL.User == nil {
		return fmt.Errorf("DATABASE_URL must be a valid PostgreSQL connection URL")
	}

	if c.Database.MaxOpenConns <= 0 ||
		c.Database.MaxIdleConns < 0 ||
		c.Database.MaxIdleConns > c.Database.MaxOpenConns {
		return fmt.Errorf("invalid database connection pool settings")
	}

	return nil
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}
