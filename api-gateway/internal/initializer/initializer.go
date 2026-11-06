package initializer

import (
	"deployflow/api-gateway/internal/config"
	"deployflow/api-gateway/internal/middleware"
	"deployflow/api-gateway/internal/routes"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func Initialize(cfg *config.Config) (*http.Server, error) {
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	router := gin.New()

	// Register global middleware
	rateLimiter := middleware.NewRateLimiter(100, time.Minute)
	router.Use(
		gin.Logger(),
		gin.Recovery(),
		middleware.RequestID(),
		middleware.CORSMiddleware(),
		middleware.SecurityHeaders(),
		rateLimiter.Middleware(),
	)

	if err := routes.Register(router, cfg); err != nil {
		return nil, err
	}

	server := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      router,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
	}

	return server, nil
}
