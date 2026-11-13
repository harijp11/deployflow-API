package dtos

import "time"

type RegisterRequestDTO struct {
	FirstName   string  `json:"first_name" binding:"required,min=2,max=100"`
	LastName    string  `json:"last_name" binding:"required,min=1,max=100"`
	Email       string  `json:"email" binding:"required,email"`
	Password    string  `json:"password" binding:"required,min=8"`
	CompanyName *string `json:"company_name"`
	PhoneNumber *string `json:"phone_number"`
}

type LoginRequestDTO struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type RefreshTokenRequestDTO struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type LogoutRequestDTO struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type ForgotPasswordRequestDTO struct {
	Email string `json:"email" binding:"required,email"`
}

type ResetPasswordRequestDTO struct {
	Token       string `json:"token" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=8"`
}

type VerifyEmailRequestDTO struct {
	Token string `json:"token" binding:"required"`
}

type TokenResponseDTO struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	TokenType    string    `json:"token_type"`
	ExpiresAt    time.Time `json:"expires_at"`
}

type RegisterResponseDTO struct {
	Message string `json:"message"`
}

type LoginResponseDTO struct {
	Message string `json:"message"`
	Tokens  TokenResponseDTO `json:"tokens"`
}

type RefreshTokenResponseDTO struct {
	Message string `json:"message"`
	Tokens  TokenResponseDTO `json:"tokens"`
}

type LogoutResponseDTO struct {
	Message string `json:"message"`
}

type ForgotPasswordResponseDTO struct {
	Message string `json:"message"`
}

type ResetPasswordResponseDTO struct {
	Message string `json:"message"`
}

type VerifyEmailResponseDTO struct {
	Message string `json:"message"`
}
