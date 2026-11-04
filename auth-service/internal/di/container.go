
package di

import (
	"log/slog"

	"deployflow/auth-service/internal/config"
	"deployflow/auth-service/internal/health"

	"gorm.io/gorm"
)

type Container struct {
	Config       *config.Config
	DB           *gorm.DB
	Logger       *slog.Logger
	HealthHandler *health.Handler
}

func NewContainer(
	cfg *config.Config,
	db *gorm.DB,
	logger *slog.Logger,
) *Container {
	return &Container{
		Config:        cfg,
		DB:            db,
		Logger:        logger,
		HealthHandler: health.NewHandler(db),
	}
}
