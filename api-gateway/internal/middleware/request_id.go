package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const RequestIDHeader = "X-Request-ID"

func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader(RequestIDHeader)

		// Generate a new ID if the client didn't provide one.
		if requestID == "" {
			requestID = uuid.NewString()
		}

		// Store the ID in the request context.
		c.Set("request_id", requestID)

		// Return the ID in the response.
		c.Header(RequestIDHeader, requestID)

		c.Next()
	}
}
