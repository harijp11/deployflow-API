package initializer

import (
	"log/slog"

	"deployflow/auth-service/internal/di"
	"deployflow/auth-service/internal/routes"
	"deployflow/shared/middleware"

	"github.com/gin-gonic/gin"
)

func InitializeApp(deps *di.Container) *gin.Engine {
	gin.SetMode(deps.Config.App.GinMode)

	router := gin.New()

	registerMiddleware(router, deps.Logger)

	routes.RegisterAuthRoutes(router, deps.AuthHandler)

	return router
}

func registerMiddleware(router *gin.Engine, logger *slog.Logger) {
	router.Use(
		gin.Recovery(),
		middleware.RequestIDMiddleware(),
		middleware.AccessLogMiddleware(),
	)
}