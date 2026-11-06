package routes

import (
	"deployflow/auth-service/internal/di"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine, deps *di.Container) {
	// router.POST("/api/v1/auth/register", deps.AuthHandler.Register)
	// router.POST("/api/v1/auth/login", deps.AuthHandler.Login)
	// router.POST("/api/v1/auth/refresh", deps.AuthHandler.RefreshToken)
	// router.POST("/api/v1/auth/logout", deps.AuthHandler.Logout)
	// router.POST("/api/v1/auth/forgot-password", deps.AuthHandler.ForgotPassword)
	// router.POST("/api/v1/auth/reset-password", deps.AuthHandler.ResetPassword)
	// router.POST("/api/v1/auth/verify-email", deps.AuthHandler.VerifyEmail)
}