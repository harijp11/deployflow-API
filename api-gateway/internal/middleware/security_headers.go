package middleware

import "github.com/gin-gonic/gin"

func SecurityHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Prevent browsers from interpreting files as a different MIME type.
		c.Header("X-Content-Type-Options", "nosniff")

		// Prevent the application from being embedded in frames.
		c.Header("X-Frame-Options", "DENY")

		// Control referrer information sent to other websites.
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")

		// Disable unnecessary browser features.
		c.Header(
			"Permissions-Policy",
			"camera=(), microphone=(), geolocation=()",
		)

		c.Next()
	}
}
