package handlers

import (
	"log"
	"net/http"

	"deployflow/user-service/internal/dtos"
	"deployflow/user-service/internal/services"

	"github.com/gin-gonic/gin"
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
        log.Printf("USER VALIDATION ERROR: %v", err)
        c.JSON(http.StatusBadRequest, gin.H{
            "message": "Invalid request",
        })
        return
    }

    registrationID := c.GetHeader("Idempotency-Key")

    if registrationID == "" {
        log.Println("USER ERROR: Idempotency-Key header is missing")
        c.JSON(http.StatusBadRequest, gin.H{
            "message": "Idempotency-Key is required",
        })
        return
    }

    userID, err := h.userService.CreateUser(
    c.Request.Context(),
    registrationID,
    req,
)
if err != nil {
    log.Printf("USER SERVICE ERROR: %v", err)

    c.JSON(http.StatusInternalServerError, gin.H{
        "message": "Failed to create user",
    })
    return
}

    c.JSON(http.StatusCreated, gin.H{
        "message": "User created successfully",
        "data": dtos.CreateUserResponseDTO{
            ID: userID,
        },
    })
}