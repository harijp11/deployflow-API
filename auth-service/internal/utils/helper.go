package utils

import "github.com/gin-gonic/gin"

func GetCurrentUserID(c *gin.Context) (uint, error) {
	userID, exists := c.Get("userID")
	if !exists {
		return 0, NewUnauthorizedError("Unauthorized access")
	}

	id, ok := userID.(uint)
	if !ok {
		return 0, NewUnauthorizedError("Invalid user ID")
	}

	return id, nil
}