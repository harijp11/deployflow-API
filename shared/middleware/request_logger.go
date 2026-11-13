
package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	reset  = "\033[0m"
	red    = "\033[31m"
	green  = "\033[32m"
	yellow = "\033[33m"
	blue   = "\033[34m"
	cyan   = "\033[36m"
	gray   = "\033[90m"
)

func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader("X-Request-ID")

		if _, err := uuid.Parse(requestID); err != nil {
			requestID = uuid.NewString()
		}

		c.Set("request_id", requestID)
		c.Header("X-Request-ID", requestID)

		c.Next()
	}
}

func AccessLogMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		duration := time.Since(start)
		status := c.Writer.Status()
		method := c.Request.Method

		fmt.Printf(
			"%s%-7s%s %s%-3d%s %s%s%s %s%v%s %s%d B%s\n",
			methodColor(method), method, reset,
			statusColor(status), status, reset,
			gray, c.Request.URL.Path, reset,
			cyan, duration, reset,
			gray, c.Writer.Size(), reset,
		)
	}
}

func methodColor(method string) string {
	switch method {
	case http.MethodGet:
		return green
	case http.MethodPost:
		return yellow
	case http.MethodPut, http.MethodPatch:
		return blue
	case http.MethodDelete:
		return red
	default:
		return cyan
	}
}

func statusColor(status int) string {
	switch {
	case status >= 200 && status < 300:
		return green
	case status >= 300 && status < 400:
		return cyan
	case status >= 400 && status < 500:
		return yellow
	case status >= 500:
		return red
	default:
		return gray
	}
}
