package services

import (
    "context"

    "user-service/internal/dtos"
    "user-service/internal/models"
    "user-service/internal/repositories"
)

type UserServiceInterface interface {
    CreateUser(ctx context.Context, req dtos.CreateUserRequestDTO) (uint, error)
}

type UserService struct {
    repo repositories.UserRepository
}

func NewUserService(repo repositories.UserRepository) UserServiceInterface {
    return &UserService{repo: repo}
}

func (us *UserService) CreateUser(
    ctx context.Context,
    req dtos.CreateUserRequestDTO,
) (uint, error) {
    user := &models.User{
        FirstName:   req.FirstName,
        LastName:    req.LastName,
        Email:       req.Email,
        CompanyName: req.CompanyName,
        PhoneNumber: req.PhoneNumber,
    }

    if err := us.repo.Create(ctx, user); err != nil {
        return 0, err
    }

    return user.ID, nil
}