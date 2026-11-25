package routes

import (
    "deployflow/user-service/internal/handlers"

    "github.com/gin-gonic/gin"
)

func RegisterUserRoutes(router *gin.Engine, userHandler *handlers.UserHandler) {
    users := router.Group("/api/v1/users")
    {
        users.POST("", userHandler.CreateUser)

        users.GET("/me", userHandler.GetMyProfile)
        users.PATCH("/me", userHandler.UpdateMyProfile)
        // users.DELETE("/me", userHandler.DeleteMyAccount)

        users.GET("", userHandler.GetAllUsers)
        users.GET("/:id", userHandler.GetUserByID)
    }
}