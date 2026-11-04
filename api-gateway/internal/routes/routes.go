package routes

import (
	"net/http"

	"deployflow/api-gateway/internal/config"
	"deployflow/api-gateway/internal/proxy"

	"github.com/gin-gonic/gin"
)

func Register(router *gin.Engine, cfg *config.Config) error {
	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"service": "api-gateway",
		})
	})

	// Register backend service routes
	services := []struct {
		name   string
		path   string
		target string
	}{
		{"auth", "/api/v1/auth", cfg.AuthServiceURL},
		{"users", "/api/v1/users", cfg.UserServiceURL},
		{"projects", "/api/v1/projects", cfg.ProjectServiceURL},
		{"deployments", "/api/v1/deployments", cfg.DeploymentServiceURL},
		{"builds", "/api/v1/builds", cfg.BuildServiceURL},
		{"runtime", "/api/v1/runtime", cfg.RuntimeServiceURL},
	}

	for _, service := range services {
		if err := proxy.RegisterProxyRoutes(
			router,
			service.name,
			service.path,
			service.target,
		); err != nil {
			return err
		}
	}

	return nil
}
