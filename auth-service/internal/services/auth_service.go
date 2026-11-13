package services

import (
	"context"
	"errors"
	"fmt"
	"net/mail"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"deployflow/auth-service/internal/config"
	"deployflow/auth-service/internal/dtos"
	"deployflow/auth-service/internal/models"
	"deployflow/auth-service/internal/repositories"
	"deployflow/auth-service/internal/utils"

	"github.com/google/uuid"
)

type AuthServiceInterface interface {
	Register(ctx context.Context, req dtos.RegisterRequestDTO) error
	Login(ctx context.Context, req dtos.LoginRequestDTO) (*dtos.TokenResponseDTO, error)
	RefreshToken(ctx context.Context, req dtos.RefreshTokenRequestDTO) (*dtos.TokenResponseDTO, error)
	Logout(ctx context.Context, req dtos.LogoutRequestDTO) error
	ForgotPassword(ctx context.Context, req dtos.ForgotPasswordRequestDTO) error
	ResetPassword(ctx context.Context, req dtos.ResetPasswordRequestDTO) error
	VerifyEmail(ctx context.Context, req dtos.VerifyEmailRequestDTO) error
}

type AuthService struct {
	repo        repositories.AuthRepository
	emailSender EmailSender
	jwtService  JWTServiceInterface
	userCreator UserCreator
	config      config.JWTConfig
}

func NewAuthService(
	repo repositories.AuthRepository,
	emailSender EmailSender,
	jwtService JWTServiceInterface,
	userCreator UserCreator,
	cfg config.JWTConfig,
) AuthServiceInterface {
	return &AuthService{
		repo:        repo,
		emailSender: emailSender,
		jwtService:  jwtService,
		userCreator: userCreator,
		config:      cfg,
	}
}

func (as *AuthService) Register(ctx context.Context, req dtos.RegisterRequestDTO) error {
	req.Email = strings.ToLower(strings.TrimSpace(req.Email))

	if err := validateEmail(req.Email); err != nil {
		return err
	}

	if err := validatePassword(req.Password); err != nil {
		return err
	}

	if as.emailSender == nil {
		return errors.New("email sender is not configured")
	}

	existingCredential, err := as.repo.FindCredentialByEmail(ctx, req.Email)
	if err != nil {
		return err
	}
	if existingCredential != nil {
		return utils.NewConflictError("Email already registered")
	}

	passwordHash, err := utils.HashPassword(req.Password)
	if err != nil {
		return fmt.Errorf("hash password: %w", err)
	}

	token, err := utils.GenerateSecureToken()
	if err != nil {
		return fmt.Errorf("generate registration token: %w", err)
	}

	now := time.Now().UTC()

	intent := &models.RegistrationIntent{
		FirstName:             req.FirstName,
		LastName:              req.LastName,
		Email:                 req.Email,
		PasswordHash:          passwordHash,
		CompanyName:           req.CompanyName,
		PhoneNumber:           req.PhoneNumber,
		VerificationTokenHash: utils.HashToken(token),
		ExpiresAt:             now.Add(as.config.EmailVerificationTTL),
		EmailSendStatus:       "pending",
		CreatedAt:             now,
		UpdatedAt:             now,
	}

	if err := as.repo.CreateOrUpdateRegistrationIntent(ctx, intent); err != nil {
		return err
	}

	if err := as.emailSender.SendVerificationEmail(ctx, req.Email, token); err != nil {
		_ = as.repo.MarkRegistrationEmailFailed(ctx, intent.ID)
		return fmt.Errorf("send verification email: %w", err)
	}

	return as.repo.MarkRegistrationEmailSent(ctx, intent.ID)
}

func (as *AuthService) Login(ctx context.Context, req dtos.LoginRequestDTO) (*dtos.TokenResponseDTO, error) {
	req.Email = strings.ToLower(strings.TrimSpace(req.Email))

	if err := validateEmail(req.Email); err != nil {
		return nil, err
	}

	credential, err := as.repo.FindCredentialByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}

	if credential == nil {
		return nil, utils.NewUnauthorizedError("Invalid email or password")
	}

	if err := utils.ComparePassword(credential.PasswordHash, req.Password); err != nil {
		return nil, utils.NewUnauthorizedError("Invalid email or password")
	}

	if credential.EmailVerifiedAt == nil {
		return nil, utils.NewForbiddenError("Email is not verified")
	}

	return as.createTokenPair(ctx, credential, uuid.New())
}

func (as *AuthService) RefreshToken(ctx context.Context, req dtos.RefreshTokenRequestDTO) (*dtos.TokenResponseDTO, error) {
	if strings.TrimSpace(req.RefreshToken) == "" {
		return nil, utils.NewUnauthorizedError("Invalid or expired token")
	}

	session, err := as.repo.FindSessionByTokenHash(ctx, utils.HashToken(req.RefreshToken))
	if err != nil {
		return nil, err
	}

	if session == nil || session.RevokedAt != nil || !time.Now().Before(session.ExpiresAt) {
		return nil, utils.NewUnauthorizedError("Invalid or expired token")
	}

	credential, err := as.repo.FindCredentialByID(ctx, session.CredentialID)
	if err != nil {
		return nil, err
	}

	if credential == nil {
		return nil, utils.NewUnauthorizedError("Invalid or expired token")
	}

	if credential.EmailVerifiedAt == nil {
		return nil, utils.NewForbiddenError("Email is not verified")
	}

	newRefreshToken, err := utils.GenerateSecureToken()
	if err != nil {
		return nil, fmt.Errorf("generate refresh token: %w", err)
	}

	now := time.Now().UTC()

	newSession := &models.AuthSession{
		CredentialID: session.CredentialID,
		UserID:       session.UserID,
		TokenHash:    utils.HashToken(newRefreshToken),
		FamilyID:     session.FamilyID,
		ExpiresAt:    now.Add(as.config.RefreshTokenTTL),
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	accessToken, accessExpiresAt, err := as.jwtService.GenerateAccessToken(credential.UserID, credential.ID)
	if err != nil {
		return nil, err
	}

	if err := as.repo.RotateSession(ctx, session.ID, newSession); err != nil {
		return nil, err
	}

	return &dtos.TokenResponseDTO{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
		TokenType:    "Bearer",
		ExpiresAt:    accessExpiresAt,
	}, nil
}

func (as *AuthService) Logout(ctx context.Context, req dtos.LogoutRequestDTO) error {
	if strings.TrimSpace(req.RefreshToken) == "" {
		return utils.NewUnauthorizedError("Invalid or expired token")
	}

	session, err := as.repo.FindSessionByTokenHash(ctx, utils.HashToken(req.RefreshToken))
	if err != nil {
		return err
	}

	if session == nil {
		return utils.NewUnauthorizedError("Invalid or expired token")
	}

	if session.RevokedAt != nil {
		return nil
	}

	return as.repo.RevokeSession(ctx, session.ID)
}

func (as *AuthService) ForgotPassword(ctx context.Context, req dtos.ForgotPasswordRequestDTO) error {
	req.Email = strings.ToLower(strings.TrimSpace(req.Email))

	if err := validateEmail(req.Email); err != nil {
		return err
	}

	credential, err := as.repo.FindCredentialByEmail(ctx, req.Email)
	if err != nil {
		return err
	}

	if credential == nil {
		return nil
	}

	if as.emailSender == nil {
		return errors.New("email sender is not configured")
	}

	token, err := utils.GenerateSecureToken()
	if err != nil {
		return fmt.Errorf("generate password reset token: %w", err)
	}

	now := time.Now().UTC()

	resetToken := &models.PasswordResetToken{
		CredentialID: credential.ID,
		TokenHash:    utils.HashToken(token),
		ExpiresAt:    now.Add(as.config.PasswordResetTokenTTL),
		CreatedAt:    now,
	}

	if err := as.repo.SavePasswordResetToken(ctx, resetToken); err != nil {
		return err
	}

	if err := as.emailSender.SendPasswordResetEmail(ctx, req.Email, token); err != nil {
		return fmt.Errorf("send password reset email: %w", err)
	}

	return nil
}

func (as *AuthService) ResetPassword(ctx context.Context, req dtos.ResetPasswordRequestDTO) error {
	if strings.TrimSpace(req.Token) == "" {
		return utils.NewUnauthorizedError("Invalid or expired token")
	}

	if err := validatePassword(req.NewPassword); err != nil {
		return err
	}

	resetToken, err := as.repo.FindPasswordResetToken(ctx, utils.HashToken(req.Token))
	if err != nil {
		return err
	}

	if resetToken == nil {
		return utils.NewUnauthorizedError("Invalid or expired token")
	}

	passwordHash, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		return fmt.Errorf("hash password: %w", err)
	}

	now := time.Now().UTC()

	credential := &models.Credential{
		ID:                resetToken.CredentialID,
		PasswordHash:      passwordHash,
		PasswordChangedAt: now,
		UpdatedAt:         now,
	}

	return as.repo.CompletePasswordReset(ctx, resetToken.ID, credential)
}

func (as *AuthService) VerifyEmail(ctx context.Context, req dtos.VerifyEmailRequestDTO) error {
	if strings.TrimSpace(req.Token) == "" {
		return utils.NewUnauthorizedError("Invalid or expired token")
	}

	if as.userCreator == nil {
		return errors.New("user service is not configured")
	}

	intent, err := as.repo.FindRegistrationIntentByTokenHash(
		ctx,
		utils.HashToken(req.Token),
	)
	if err != nil {
		return err
	}

	if intent == nil {
		return utils.NewUnauthorizedError("Invalid or expired token")
	}

	profile := UserProfile{
		FirstName:   intent.FirstName,
		LastName:    intent.LastName,
		Email:       intent.Email,
		CompanyName: intent.CompanyName,
		PhoneNumber: intent.PhoneNumber,
	}

	registrationID := strconv.FormatUint(uint64(intent.ID), 10)
	userID, err := as.userCreator.CreateUser(ctx, registrationID, profile)
	if err != nil {
		return fmt.Errorf("create user profile: %w", err)
	}

	now := time.Now().UTC()

	credential := &models.Credential{
		UserID:            userID,
		Email:             intent.Email,
		PasswordHash:      intent.PasswordHash,
		EmailVerifiedAt:   &now,
		PasswordChangedAt: now,
		CreatedAt:         now,
		UpdatedAt:         now,
	}

	if err := as.repo.CompleteRegistration(ctx, intent.ID, credential); err != nil {
		return fmt.Errorf("complete registration: %w", err)
	}

	return nil
}

func (as *AuthService) createTokenPair(ctx context.Context, credential *models.Credential, familyID uuid.UUID) (*dtos.TokenResponseDTO, error) {
	accessToken, accessExpiresAt, err := as.jwtService.GenerateAccessToken(credential.UserID, credential.ID)
	if err != nil {
		return nil, err
	}

	refreshToken, err := utils.GenerateSecureToken()
	if err != nil {
		return nil, fmt.Errorf("generate refresh token: %w", err)
	}

	now := time.Now().UTC()

	session := &models.AuthSession{
		CredentialID: credential.ID,
		UserID:       credential.UserID,
		TokenHash:    utils.HashToken(refreshToken),
		FamilyID:     familyID,
		ExpiresAt:    now.Add(as.config.RefreshTokenTTL),
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	if err := as.repo.CreateSession(ctx, session); err != nil {
		return nil, err
	}

	return &dtos.TokenResponseDTO{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresAt:    accessExpiresAt,
	}, nil
}

func validateEmail(email string) error {
	parsedEmail, err := mail.ParseAddress(email)
	if err != nil || parsedEmail.Address != email {
		return utils.NewValidationError("Invalid email address", "INVALID_INPUT", "email")
	}

	return nil
}

func validatePassword(password string) error {
	if utf8.RuneCountInString(password) < 8 {
		return utils.NewValidationError("Password must be at least 8 characters", "INVALID_INPUT", "password")
	}

	return nil
}
