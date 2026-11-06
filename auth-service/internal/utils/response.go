package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Response(c *gin.Context, statusCode int, message string, data any) {
	switch statusCode {
	case http.StatusOK, http.StatusCreated:
		response := gin.H{
			"message": message,
		}

		if data != nil {
			response["data"] = data
		}

		c.JSON(statusCode, response)

	case http.StatusBadRequest,
		http.StatusUnauthorized,
		http.StatusForbidden,
		http.StatusNotFound,
		http.StatusConflict:
		c.JSON(statusCode, gin.H{
			"error": message,
		})

	default:
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "An unexpected error occurred",
		})
	}
}

func SendNotFoundErrResp(c *gin.Context, err *ResourceNotFoundError) {
	c.JSON(http.StatusNotFound, gin.H{
		"error": err.Message,
	})
}

func SendInvalidErrResp(c *gin.Context, err *ValidationError) {
	c.JSON(http.StatusBadRequest, gin.H{
		"error": err.Message,
		"code":  err.Code,
		"field": err.Field,
	})
}

func SendConflictErrResp(c *gin.Context, err *ConflictError) {
	c.JSON(http.StatusConflict, gin.H{
		"error": err.Message,
	})
}

func SendUnauthorizedErrResp(c *gin.Context, err *UnauthorizedError) {
	c.JSON(http.StatusUnauthorized, gin.H{
		"error": err.Message,
	})
}

func SendForbiddenErrResp(c *gin.Context, err *ForbiddenError) {
	c.JSON(http.StatusForbidden, gin.H{
		"error": err.Message,
	})
}

func SendBadRequestErrResp(c *gin.Context, err *BadRequestError) {
	c.JSON(http.StatusBadRequest, gin.H{
		"error": err.Message,
	})
}

func SendInternalErrResp(c *gin.Context) {
	c.JSON(http.StatusInternalServerError, gin.H{
		"error": "An unexpected error occurred",
	})
}