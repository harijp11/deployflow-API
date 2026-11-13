package di

import (
    "log/slog"

    "deployflow/user-service/internal/config"
    "deployflow/user-service/internal/handlers"
    "deployflow/user-service/internal/repositories"
    "deployflow/user-service/internal/services"

    "gorm.io/gorm"
)

type Container struct {
    Config      *config.Config
    DB          *gorm.DB
    Logger      *slog.Logger
    UserHandler *handlers.UserHandler
}

func NewContainer(cfg *config.Config, db *gorm.DB, logger *slog.Logger) *Container {
    // Repositories
    userRepo := repositories.NewUserRepository(db)

    // Services
    userService := services.NewUserService(userRepo)

    // Handlers
    userHandler := handlers.NewUserHandler(userService)

    return &Container{
        Config:      cfg,
        DB:          db,
        Logger:      logger,
        UserHandler: userHandler,
    }
}