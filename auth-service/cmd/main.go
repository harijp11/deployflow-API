
package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"deployflow/auth-service/internal/config"
	"deployflow/auth-service/internal/database"
	"deployflow/auth-service/internal/di"
	"deployflow/auth-service/internal/initializer"
)

func main() {
	if err := run(); err != nil {
		log.Fatalf("Auth Service failed: %v", err)
	}
}

func run() error {
	// 1. Load configuration.
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	// 2. Initialize structured logger.
	logger := slog.New(
		slog.NewJSONHandler(os.Stdout, nil),
	)

	// 3. Connect to PostgreSQL.
	db, err := database.Connect(cfg.Database)
	if err != nil {
		return fmt.Errorf("connect database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("get database connection: %w", err)
	}
	defer sqlDB.Close()

	// 4. Initialize dependency container.
	deps := di.NewContainer(cfg, db, logger)

	// 5. Initialize HTTP application.
	router := initializer.InitializeApp(deps)

	// 6. Configure HTTP server.
	server := &http.Server{
		Addr:              ":" + cfg.Server.Port,
		Handler:           router,
		ReadTimeout:       cfg.Server.ReadTimeout,
		WriteTimeout:      cfg.Server.WriteTimeout,
		IdleTimeout:       cfg.Server.IdleTimeout,
		ReadHeaderTimeout: cfg.Server.ReadHeaderTimeout,
	}

	// 7. Start HTTP server.
	serverErrors := make(chan error, 1)

	go func() {
		logger.Info("Auth Service starting",
			"port", cfg.Server.Port,
			"environment", cfg.App.Environment,
		)

		serverErrors <- server.ListenAndServe()
	}()

	// 8. Wait for shutdown signal or server failure.
	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	select {
	case err := <-serverErrors:
		if !errors.Is(err, http.ErrServerClosed) {
			return fmt.Errorf("HTTP server: %w", err)
		}

	case <-ctx.Done():
		logger.Info("Shutdown signal received")
	}

	// 9. Graceful shutdown.
	shutdownCtx, cancel := context.WithTimeout(
		context.Background(),
		cfg.Server.ShutdownTimeout,
	)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		_ = server.Close()
		return fmt.Errorf("graceful shutdown: %w", err)
	}

	logger.Info("Auth Service stopped successfully")

	return nil
}
