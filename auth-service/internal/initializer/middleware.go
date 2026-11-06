package initializer

import (
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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

		if _, err := uuid.Parse(requestID); err != nil {
			requestID = generateRequestID()
		}

		c.Set("request_id", requestID)
		c.Header("X-Request-ID", requestID)

		c.Next()
	}
}

func accessLogMiddleware(logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		logger.Info("HTTP request",
			"request_id", c.GetString("request_id"),
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"status", c.Writer.Status(),
			"latency_ms", time.Since(start).Milliseconds(),
			"response_size", c.Writer.Size(),
			"client_ip", c.ClientIP(),
			"protocol", c.Request.Proto,
		)
	}
}

func generateRequestID() string {
	return uuid.NewString()
}