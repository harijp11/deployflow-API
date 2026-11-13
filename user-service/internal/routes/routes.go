package routes

import (
    "deployflow/user-service/internal/handlers"

    "github.com/gin-gonic/gin"
)

func RegisterUserRoutes(router *gin.Engine, userHandler *handlers.UserHandler) {
    users := router.Group("/api/v1/users")
    {
        users.POST("", userHandler.CreateUser)
    }
}