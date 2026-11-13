package routes

import (
    "github.com/gin-gonic/gin"
    "user-service/internal/handlers"
)

func RegisterUserRoutes(router *gin.Engine, userHandler *handlers.UserHandler) {
    users := router.Group("/api/v1/users")
    {
        users.POST("", userHandler.CreateUser)
    }
}