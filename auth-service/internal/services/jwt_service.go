package services

import (
	"fmt"
	"time"

	"deployflow/auth-service/internal/config"

	"github.com/golang-jwt/jwt/v5"
)

type JWTServiceInterface interface {
	GenerateAccessToken(userID, credentialID uint) (string, time.Time, error)
}

type JWTService struct {
	config config.JWTConfig
}

func NewJWTService(config config.JWTConfig) JWTServiceInterface {
	return &JWTService{
		config: config,
	}
}

func (js *JWTService) GenerateAccessToken(
	userID, credentialID uint,
) (string, time.Time, error) {
	now := time.Now()
	expiresAt := now.Add(js.config.AccessTokenTTL)

	claims := jwt.MapClaims{
		"iss":           js.config.Issuer,
		"sub":           fmt.Sprintf("%d", userID),
		"aud":           jwt.ClaimStrings{js.config.Audience},
		"iat":           now.Unix(),
		"exp":           expiresAt.Unix(),
		"user_id":       userID,
		"credential_id": credentialID,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString([]byte(js.config.SigningKey))
	if err != nil {
		return "", time.Time{}, fmt.Errorf("sign access token: %w", err)
	}

	return signedToken, expiresAt, nil
}