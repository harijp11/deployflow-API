package handler

import "github.com/gin-gonic/gin"

type AuthHandler struct {
    // Dependencies will be injected here.
}

func (h *AuthHandler) Register(c *gin.Context) {
}

func (h *AuthHandler) Login(c *gin.Context) {
}

func (h *AuthHandler) RefreshToken(c *gin.Context) {
}

func (h *AuthHandler) Logout(c *gin.Context) {
}

func (h *AuthHandler) ForgotPassword(c *gin.Context) {
}

func (h *AuthHandler) ResetPassword(c *gin.Context) {
}

func (h *AuthHandler) VerifyEmail(c *gin.Context) {
}