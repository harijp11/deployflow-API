
package initializer

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

func registerMiddleware(router *gin.Engine, depsLogger *slog.Logger) {
	router.Use(
		gin.Recovery(),
		requestIDMiddleware(),
		accessLogMiddleware(depsLogger),
	)
}

func requestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader("X-Request-ID")

		if requestID == "" {
			requestID = generateRequestID()
		}

		c.Set("request_id", requestID)
		c.Header("X-Request-ID", requestID)

		c.Next()
	}
}

func accessLogMiddleware(logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		logger.Info("HTTP request",
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"status", c.Writer.Status(),
			"latency", c.Writer.Size(),
			"client_ip", c.ClientIP(),
		)
	}
}

func generateRequestID() string {
	// Replace with a UUID generator.
	return http.CanonicalHeaderKey("X-Request-ID")
}
