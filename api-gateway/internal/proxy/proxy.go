package proxy

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
)

func NewReverseProxy(target string) (gin.HandlerFunc, error) {
	targetURL, err := url.Parse(target)
	if err != nil {
		return nil, fmt.Errorf("invalid proxy target URL: %w", err)
	}

	if (targetURL.Scheme != "http" && targetURL.Scheme != "https") ||
		targetURL.Host == "" ||
		targetURL.User != nil {
		return nil, fmt.Errorf("proxy target must be a valid HTTP or HTTPS URL")
	}

	reverseProxy := httputil.NewSingleHostReverseProxy(targetURL)

	reverseProxy.ErrorHandler = func(
		w http.ResponseWriter,
		r *http.Request,
		err error,
	) {
		http.Error(w, "Backend service unavailable", http.StatusBadGateway)
	}

	return func(c *gin.Context) {
		// Retrieve the request ID from Gin's context.
		if requestID, exists := c.Get("request_id"); exists {
			if id, ok := requestID.(string); ok && id != "" {
				c.Request.Header.Set("X-Request-ID", id)
			}
		}

		reverseProxy.ServeHTTP(c.Writer, c.Request)
	}, nil
}

func RegisterProxyRoutes(router *gin.Engine, serviceName string, basePath string, target string) error {
	handler, err := NewReverseProxy(target)
	if err != nil {
		return fmt.Errorf("%s proxy: %w", serviceName, err)
	}

	basePath = strings.TrimRight(basePath, "/")

	router.Any(basePath, handler)
	router.Any(basePath+"/*path", handler)

	return nil
}