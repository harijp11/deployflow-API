package handlers

import (
    "net/http"

    "github.com/gin-gonic/gin"
    "user-service/internal/dtos"
    "user-service/internal/services"
)

type UserHandler struct {
    userService services.UserServiceInterface
}

func NewUserHandler(userService services.UserServiceInterface) *UserHandler {
    return &UserHandler{userService: userService}
}

func (h *UserHandler) CreateUser(c *gin.Context) {
    var req dtos.CreateUserRequestDTO

    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request"})
        return
    }

    userID, err := h.userService.CreateUser(c.Request.Context(), req)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create user"})
        return
    }

    c.JSON(http.StatusCreated, gin.H{
        "message": "User created successfully",
        "data": dtos.CreateUserResponseDTO{
            ID: userID,
        },
    })
}