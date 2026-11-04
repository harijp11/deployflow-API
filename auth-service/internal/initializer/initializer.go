
package initializer

import (
	"deployflow/auth-service/internal/di"
	"deployflow/auth-service/internal/health"
	"deployflow/auth-service/internal/routes"

	"github.com/gin-gonic/gin"
)

func InitializeApp(deps *di.Container) *gin.Engine {
	gin.SetMode(deps.Config.App.GinMode)

	router := gin.New()

	// Register global middleware.
	registerMiddleware(router, deps.Logger)

	// Health and readiness endpoints.
	healthHandler := health.NewHandler(deps.DB)
	health.RegisterRoutes(router, healthHandler)

	// Application API routes.
	routes.RegisterRoutes(router, deps)

	return router
}
