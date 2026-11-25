package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"deployflow/shared/utils"
	"deployflow/user-service/internal/dtos"
	"deployflow/user-service/internal/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserHandler struct {
	userService services.UserServiceInterface
}

func NewUserHandler(userService services.UserServiceInterface) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// POST /api/v1/users
func (h *UserHandler) CreateUser(c *gin.Context) {
	var req dtos.CreateUserRequestDTO

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Response(c, http.StatusBadRequest, "Invalid request body", nil)
		return
	}

	registrationID := c.GetHeader("Idempotency-Key")
	if registrationID == "" {
		utils.RespondAuthError(c, utils.NewBadRequestError("Idempotency-Key is required"))
		return
	}

	userID, err := h.userService.CreateUser(
		c.Request.Context(),
		registrationID,
		req,
	)
	if err != nil {
		utils.RespondAuthError(c, err)
		return
	}

	utils.Response(c, http.StatusCreated, "User created successfully",
		dtos.CreateUserResponseDTO{ID: userID},
	)
}

// GET /api/v1/users/me
func (h *UserHandler) GetMyProfile(c *gin.Context) {
	userID := c.GetUint("user_id")

	if userID == 0 {
		utils.Response(c, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}

	h.getUserProfile(c, userID)
}

// GET /api/v1/users/:id
func (h *UserHandler) GetUserByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		utils.Response(c, http.StatusBadRequest, "Invalid user ID", nil)
		return
	}

	h.getUserProfile(c, uint(id))
}

func (h *UserHandler) getUserProfile(c *gin.Context, userID uint) {
	user, err := h.userService.GetUserByID(c.Request.Context(), userID)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		utils.Response(c, http.StatusNotFound, "User not found", nil)
		return
	}

	if err != nil {
		utils.RespondAuthError(c, err)
		return
	}

	utils.Response(c, http.StatusOK, "User retrieved successfully", user)
}

// PATCH /api/v1/users/me
func (h *UserHandler) UpdateMyProfile(c *gin.Context) {
	userID := c.GetUint("user_id")

	if userID == 0 {
		utils.Response(c, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}

	var req dtos.UpdateUserRequestDTO

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Response(c, http.StatusBadRequest, "Invalid request body", nil)
		return
	}

	err := h.userService.UpdateUser(c.Request.Context(), userID, req)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		utils.Response(c, http.StatusNotFound, "User not found", nil)
		return
	}

	if err != nil {
		utils.RespondAuthError(c, err)
		return
	}

	h.getUserProfile(c, userID)
}


func (h *UserHandler) GetAllUsers(c *gin.Context) {
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil {
		utils.Response(c, http.StatusBadRequest, "Invalid page number", nil)
		return
	}

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil {
		utils.Response(c, http.StatusBadRequest, "Invalid limit. It must be between 1 and 100", nil)
		return
	}

	users, total, err := h.userService.GetAllUsers(
		c.Request.Context(),
		page,
		limit,
	)

	if err != nil {
		utils.RespondAuthError(c, err)
		return
	}

	utils.Response(c, http.StatusOK, "Users retrieved successfully", gin.H{
		"users": users,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}
