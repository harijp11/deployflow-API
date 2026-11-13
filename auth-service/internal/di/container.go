package di

import (
    "log/slog"

    "deployflow/auth-service/internal/clients"
    "deployflow/auth-service/internal/config"
    "deployflow/auth-service/internal/handlers"
    "deployflow/auth-service/internal/health"
    "deployflow/auth-service/internal/repositories"
    "deployflow/auth-service/internal/services"
    "deployflow/shared/email"

    "gorm.io/gorm"
)

type Container struct {
    Config        *config.Config
    DB            *gorm.DB
    Logger        *slog.Logger
    HealthHandler *health.Handler
    AuthHandler   *handlers.AuthHandler
}

func NewContainer(cfg *config.Config, db *gorm.DB, logger *slog.Logger) *Container {
    // External clients
    userServiceClient := clients.NewUserServiceClient(cfg.UserService.BaseURL)

    // Third-party services
    jwtService := services.NewJWTService(cfg.JWT)
    emailSender := email.NewSMTPSender(
        cfg.Email.Host,
        cfg.Email.Port,
        cfg.Email.Username,
        cfg.Email.Password,
        cfg.Email.From,
    )
    authEmailSender := services.NewAuthEmailSender(emailSender)

    // Repositories
    authRepo := repositories.NewAuthRepository(db)

    // Services
    authService := services.NewAuthService(
        authRepo,
        authEmailSender,
        jwtService,
        userServiceClient,
        cfg.JWT,
		cfg.Registration,
    )

    // Handlers
    authHandler := handlers.NewAuthHandler(authService)

    return &Container{
        Config:        cfg,
        DB:            db,
        Logger:        logger,
        HealthHandler: health.NewHandler(db),
        AuthHandler:   authHandler,
    }
}