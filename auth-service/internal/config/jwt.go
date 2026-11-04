
package config

import (
	"fmt"
	"time"
)

type JWTConfig struct {
	Issuer          string
	Audience        string
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
}

func loadJWTConfig() (*JWTConfig, error) {
	cfg := &JWTConfig{
		Issuer:          getEnv("JWT_ISSUER", "deployflow-auth"),
		Audience:        getEnv("JWT_AUDIENCE", "deployflow-api"),
		AccessTokenTTL:  15 * time.Minute,
		RefreshTokenTTL: 7 * 24 * time.Hour,
	}

	if cfg.Issuer == "" || cfg.Audience == "" {
		return nil, fmt.Errorf("JWT issuer and audience are required")
	}

	accessTTL, err := time.ParseDuration(
		getEnv("ACCESS_TOKEN_TTL", "15m"),
	)
	if err != nil || accessTTL <= 0 {
		return nil, fmt.Errorf("invalid ACCESS_TOKEN_TTL")
	}

	refreshTTL, err := time.ParseDuration(
		getEnv("REFRESH_TOKEN_TTL", "168h"),
	)
	if err != nil || refreshTTL <= accessTTL {
		return nil, fmt.Errorf("REFRESH_TOKEN_TTL must exceed ACCESS_TOKEN_TTL")
	}

	cfg.AccessTokenTTL = accessTTL
	cfg.RefreshTokenTTL = refreshTTL

	return cfg, nil
}
