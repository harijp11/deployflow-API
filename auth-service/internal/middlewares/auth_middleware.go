
package middleware

import (
	"math"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware(signingKey, issuer, audience string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. Get Authorization header
		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			unauthorized(c, "Authorization header is required")
			return
		}

		// 2. Extract Bearer token
		parts := strings.Fields(authHeader)

		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			unauthorized(c, "Invalid authorization header")
			return
		}

		tokenString := parts[1]

		// 3. Parse and validate JWT
		claims := jwt.MapClaims{}

		token, err := jwt.ParseWithClaims(
			tokenString,
			claims,
			func(token *jwt.Token) (interface{}, error) {
				// Accept only HS256
				if token.Method != jwt.SigningMethodHS256 {
					return nil, jwt.ErrSignatureInvalid
				}

				return []byte(signingKey), nil
			},
			jwt.WithValidMethods([]string{"HS256"}),
			jwt.WithIssuer(issuer),
			jwt.WithAudience(audience),
		)

		if err != nil || token == nil || !token.Valid {
			unauthorized(c, "Invalid or expired token")
			return
		}

		// 4. Extract session_id
		sessionID, ok := claims["session_id"].(float64)

		if !ok || sessionID <= 0 || math.Trunc(sessionID) != sessionID {
			unauthorized(c, "Invalid session")
			return
		}

		// 5. Extract user_id
		userID, ok := claims["user_id"].(float64)

		if !ok || userID <= 0 || math.Trunc(userID) != userID {
			unauthorized(c, "Invalid user")
			return
		}

		// 6. Extract credential_id
		credentialID, ok := claims["credential_id"].(float64)

		if !ok || credentialID <= 0 || math.Trunc(credentialID) != credentialID {
			unauthorized(c, "Invalid credential")
			return
		}

		// 7. Store values in Gin context
		c.Set("session_id", uint(sessionID))
		c.Set("user_id", uint(userID))
		c.Set("credential_id", uint(credentialID))

		// Continue to handler
		c.Next()
	}
}

func unauthorized(c *gin.Context, message string) {
	c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
		"message": message,
	})
}
