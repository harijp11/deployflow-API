package config

import (
    "fmt"
    "time"
)

type RegistrationConfig struct {
    EmailVerificationTTL time.Duration
}

func loadRegistrationConfig() (*RegistrationConfig, error) {
    ttl, err := time.ParseDuration(
        getEnv("EMAIL_VERIFICATION_TTL", "30m"),
    )
    if err != nil || ttl <= 0 {
        return nil, fmt.Errorf("invalid EMAIL_VERIFICATION_TTL")
    }

    return &RegistrationConfig{
        EmailVerificationTTL: ttl,
    }, nil
}