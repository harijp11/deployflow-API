
package config

import (
	"fmt"
	"time"
)

type JWTConfig struct {
	Issuer                string
	Audience              string
	SigningKey            string
	AccessTokenTTL        time.Duration
	RefreshTokenTTL       time.Duration
	PasswordResetTokenTTL time.Duration
	EmailVerificationTTL  time.Duration
}


func loadJWTConfig() (*JWTConfig, error) {
	cfg := &JWTConfig{
		Issuer:   getEnv("JWT_ISSUER", "deployflow-auth"),
		Audience: getEnv("JWT_AUDIENCE", "deployflow-api"),
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

	passwordResetTTL, err := time.ParseDuration(
		getEnv("PASSWORD_RESET_TOKEN_TTL", "15m"),
	)
	if err != nil || passwordResetTTL <= 0 {
		return nil, fmt.Errorf("invalid PASSWORD_RESET_TOKEN_TTL")
	}

	emailVerificationTTL, err := time.ParseDuration(
		getEnv("EMAIL_VERIFICATION_TOKEN_TTL", "24h"),
	)
	if err != nil || emailVerificationTTL <= 0 {
		return nil, fmt.Errorf("invalid EMAIL_VERIFICATION_TOKEN_TTL")
	}

	cfg.AccessTokenTTL = accessTTL
	cfg.RefreshTokenTTL = refreshTTL
	cfg.PasswordResetTokenTTL = passwordResetTTL
	cfg.EmailVerificationTTL = emailVerificationTTL

	return cfg, nil
}

