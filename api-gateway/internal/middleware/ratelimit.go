package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type client struct {
	requests  int
	windowEnd time.Time
}

type RateLimiter struct {
	mu      sync.Mutex
	clients map[string]*client
	limit   int
	window  time.Duration
}

func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	return &RateLimiter{
		clients: make(map[string]*client),
		limit:   limit,
		window:  window,
	}
}

func (rl *RateLimiter) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		now := time.Now()

		rl.mu.Lock()

		entry, exists := rl.clients[ip]

		if !exists || now.After(entry.windowEnd) {
			entry = &client{
				requests:  0,
				windowEnd: now.Add(rl.window),
			}
			rl.clients[ip] = entry
		}

		if entry.requests >= rl.limit {
			rl.mu.Unlock()

			c.Header("Retry-After", time.Until(entry.windowEnd).Round(time.Second).String())
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": "rate limit exceeded",
			})
			return
		}

		entry.requests++

		rl.mu.Unlock()

		c.Next()
	}
}
