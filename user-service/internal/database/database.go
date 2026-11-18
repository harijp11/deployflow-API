package database

import (
    "fmt"

    "deployflow/user-service/internal/config"

    "gorm.io/driver/postgres"
    "gorm.io/gorm"
)

func Connect(cfg config.DatabaseConfig) (*gorm.DB, error) {
    db, err := gorm.Open(postgres.Open(cfg.URL), &gorm.Config{})
    if err != nil {
        return nil, fmt.Errorf("connect to user database: %w", err)
    }

    sqlDB, err := db.DB()
    if err != nil {
        return nil, fmt.Errorf("get underlying database: %w", err)
    }

    if err := sqlDB.Ping(); err != nil {
        return nil, fmt.Errorf("ping user database: %w", err)
    }

    return db, nil
}