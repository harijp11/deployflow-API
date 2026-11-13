package main

import (
    "log"
    "log/slog"
    "net/http"

    "deployflow/user-service/internal/config"
    "deployflow/user-service/internal/database"
    "deployflow/user-service/internal/di"
    "deployflow/user-service/internal/initializer"
)

func main() {
    cfg, err := config.Load()
    if err != nil {
        log.Fatalf("load config: %v", err)
    }

    logger := slog.Default()

    db, err := database.Connect(cfg.Database)
    if err != nil {
        log.Fatalf("connect database: %v", err)
    }

    deps := di.NewContainer(cfg, db, logger)

    router := initializer.InitializeApp(deps)

    server := &http.Server{
        Addr:    ":" + cfg.Server.Port,
        Handler: router,
    }

    logger.Info("User Service starting", "port", cfg.Server.Port)

    if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
        log.Fatalf("start User Service: %v", err)
    }
}