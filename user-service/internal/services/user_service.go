package services

import (
	"context"
	"strings"

	"deployflow/user-service/internal/dtos"
	"deployflow/user-service/internal/models"
	"deployflow/user-service/internal/repositories"
)

type UserServiceInterface interface {
	GetAllUsers(ctx context.Context, page, limit int) ([]dtos.UserResponseDTO, int64, error)
	GetUserByID(ctx context.Context, userID uint) (*dtos.UserResponseDTO, error)
	CreateUser(ctx context.Context, registrationID string, req dtos.CreateUserRequestDTO) (uint, error)
	UpdateUser(ctx context.Context, userID uint, req dtos.UpdateUserRequestDTO) error
}

type UserService struct {
	repo repositories.UserRepository
}

func NewUserService(repo repositories.UserRepository) UserServiceInterface {
	return &UserService{repo: repo}
}

func (us *UserService) GetAllUsers(ctx context.Context, page, limit int) ([]dtos.UserResponseDTO, int64, error) {
	if page < 1 {
		page = 1
	}

	if limit < 1 || limit > 100 {
		limit = 10
	}

	offset := (page - 1) * limit

	users, total, err := us.repo.FindAll(ctx, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	responses := make([]dtos.UserResponseDTO, 0, len(users))

	for i := range users {
		responses = append(responses, *toUserResponse(&users[i]))
	}

	return responses, total, nil
}

func (us *UserService) GetUserByID(ctx context.Context, userID uint) (*dtos.UserResponseDTO, error) {
	user, err := us.repo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return toUserResponse(user), nil
}

func (us *UserService) CreateUser(ctx context.Context, registrationID string, req dtos.CreateUserRequestDTO) (uint, error) {
	user := &models.User{
		FirstName:   strings.TrimSpace(req.FirstName),
		LastName:    strings.TrimSpace(req.LastName),
		Email:       strings.ToLower(strings.TrimSpace(req.Email)),
		CompanyName: trimOptionalString(req.CompanyName),
		PhoneNumber: trimOptionalString(req.PhoneNumber),
	}

	return us.repo.CreateOrGetByRegistrationID(ctx, registrationID, user)
}

func (us *UserService) UpdateUser(ctx context.Context, userID uint, req dtos.UpdateUserRequestDTO) error {
	updates := make(map[string]interface{})

	if req.FirstName != nil {
		updates["first_name"] = strings.TrimSpace(*req.FirstName)
	}

	if req.LastName != nil {
		updates["last_name"] = strings.TrimSpace(*req.LastName)
	}

	if req.CompanyName != nil {
		updates["company_name"] = trimOptionalString(req.CompanyName)
	}

	if req.PhoneNumber != nil {
		updates["phone_number"] = trimOptionalString(req.PhoneNumber)
	}

	if len(updates) == 0 {
		return nil
	}

	return us.repo.Update(ctx, userID, updates)
}

func toUserResponse(user *models.User) *dtos.UserResponseDTO {
	return &dtos.UserResponseDTO{
		ID:          user.ID,
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		Email:       user.Email,
		CompanyName: user.CompanyName,
		PhoneNumber: user.PhoneNumber,
	}
}

func trimOptionalString(value *string) *string {
	if value == nil {
		return nil
	}

	trimmed := strings.TrimSpace(*value)
	if trimmed == "" {
		return nil
	}

	return &trimmed
}
