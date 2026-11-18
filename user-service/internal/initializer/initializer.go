package initializer

import (
	"log/slog"

	"deployflow/shared/middleware"
	"deployflow/user-service/internal/di"
	"deployflow/user-service/internal/routes"

	"github.com/gin-gonic/gin"
)

func InitializeApp(deps *di.Container) *gin.Engine {
	gin.SetMode(deps.Config.App.GinMode)

	router := gin.New()

	registerMiddleware(router, deps.Logger)

	routes.RegisterUserRoutes(router, deps.UserHandler)

	return router
}

func registerMiddleware(router *gin.Engine, logger *slog.Logger) {
	router.Use(
		gin.Recovery(),
		middleware.RequestIDMiddleware(),
		middleware.AccessLogMiddleware(),
	)
}