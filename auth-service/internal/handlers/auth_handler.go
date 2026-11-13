package handlers

import (
	"net/http"

	"deployflow/auth-service/internal/dtos"
	"deployflow/auth-service/internal/services"
	"deployflow/auth-service/internal/utils"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService services.AuthServiceInterface
}

func NewAuthHandler(authService services.AuthServiceInterface) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

func (ah *AuthHandler) Register(c *gin.Context) {
	var req dtos.RegisterRequestDTO

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Response(c, http.StatusBadRequest, "Invalid request body", nil)
		return
	}

	if err := ah.authService.Register(c.Request.Context(), req); err != nil {
		utils.RespondAuthError(c, err)
		return
	}

	utils.Response(c, http.StatusCreated, "Registration successful. Please verify your email.", nil)
}

func (ah *AuthHandler) Login(c *gin.Context) {
	var req dtos.LoginRequestDTO

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Response(c, http.StatusBadRequest, "Invalid request body", nil)
		return
	}

	tokens, err := ah.authService.Login(c.Request.Context(), req)
	if err != nil {
		utils.RespondAuthError(c, err)
		return
	}

	utils.Response(c, http.StatusOK, "Login successful", tokens)
}

func (ah *AuthHandler) RefreshToken(c *gin.Context) {
	var req dtos.RefreshTokenRequestDTO

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Response(c, http.StatusBadRequest, "Invalid request body", nil)
		return
	}

	tokens, err := ah.authService.RefreshToken(c.Request.Context(), req)
	if err != nil {
		utils.RespondAuthError(c, err)
		return
	}

	utils.Response(c, http.StatusOK, "Token refreshed successfully", tokens)
}

func (ah *AuthHandler) Logout(c *gin.Context) {
	var req dtos.LogoutRequestDTO

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Response(c, http.StatusBadRequest, "Invalid request body", nil)
		return
	}

	if err := ah.authService.Logout(c.Request.Context(), req); err != nil {
		utils.RespondAuthError(c, err)
		return
	}

	utils.Response(c, http.StatusOK, "Logout successful", nil)
}

func (ah *AuthHandler) ForgotPassword(c *gin.Context) {
	var req dtos.ForgotPasswordRequestDTO

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Response(c, http.StatusBadRequest, "Invalid request body", nil)
		return
	}

	if err := ah.authService.ForgotPassword(c.Request.Context(), req); err != nil {
		utils.RespondAuthError(c, err)
		return
	}

	utils.Response(c, http.StatusOK, "If the email is registered, a password reset link will be sent", nil)
}

func (ah *AuthHandler) ResetPassword(c *gin.Context) {
	var req dtos.ResetPasswordRequestDTO

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Response(c, http.StatusBadRequest, "Invalid request body", nil)
		return
	}

	if err := ah.authService.ResetPassword(c.Request.Context(), req); err != nil {
		utils.RespondAuthError(c, err)
		return
	}

	utils.Response(c, http.StatusOK, "Password reset successful", nil)
}

func (ah *AuthHandler) VerifyEmail(c *gin.Context) {
	var req dtos.VerifyEmailRequestDTO

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Response(c, http.StatusBadRequest, "Invalid request body", nil)
		return
	}

	if err := ah.authService.VerifyEmail(c.Request.Context(), req); err != nil {
		utils.RespondAuthError(c, err)
		return
	}

	utils.Response(c, http.StatusOK, "Email verified successfully", nil)
}
