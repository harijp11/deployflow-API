
package health

import "github.com/gin-gonic/gin"

func RegisterHealthRoutes(router *gin.Engine, handler *Handler) {
	router.GET("/health", handler.Liveness)
	router.GET("/ready", handler.Readiness)
}
