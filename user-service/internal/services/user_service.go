package services

import (
	"context"
	"strings"

	"deployflow/user-service/internal/dtos"
	"deployflow/user-service/internal/models"
	"deployflow/user-service/internal/repositories"
)

type UserServiceInterface interface {
    CreateUser(ctx context.Context, registrationID string, req dtos.CreateUserRequestDTO) (uint, error)
}

type UserService struct {
    repo repositories.UserRepository
}

func NewUserService(repo repositories.UserRepository) UserServiceInterface {
    return &UserService{repo: repo}
}

func (us *UserService) CreateUser(ctx context.Context, registrationID string, req dtos.CreateUserRequestDTO) (uint, error) {
    user := &models.User{
        FirstName:   req.FirstName,
        LastName:    req.LastName,
        Email:       strings.ToLower(strings.TrimSpace(req.Email)),
        CompanyName: req.CompanyName,
        PhoneNumber: req.PhoneNumber,
    }

    userID, err := us.repo.CreateOrGetByRegistrationID(
        ctx,
        registrationID,
        user,
    )
    if err != nil {
        return 0, err
    }

    return userID, nil
}