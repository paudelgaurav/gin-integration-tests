package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	apierrors "github.com/paudelgaurav/gin-boilerplate/pkg/api_errors"
	"github.com/paudelgaurav/gin-boilerplate/pkg/framework"
	"github.com/paudelgaurav/gin-boilerplate/pkg/utils"
	"golang.org/x/time/rate"
)

// LimiterConfig holds the rate limit configuration for clients
type LimiterConfig struct {
	RequestsPerSecond rate.Limit
	BurstSize         int
	CleanupInterval   time.Duration
}

// clientLimiter contains the rate limiter and the last acess time
type clientLimiter struct {
	limiter    *rate.Limiter
	lastAccess time.Time
}

// NewRateLimitMiddleware returns a rate-limiting middleware
func NewRateLimitMiddleware(cfg LimiterConfig, logger framework.Logger) gin.HandlerFunc {

	clients := make(map[string]*clientLimiter)
	mu := sync.Mutex{}

	go cleanupClients(clients, &mu, cfg.CleanupInterval)

	return func(c *gin.Context) {
		clientIP := c.ClientIP()

		mu.Lock()

		limiter, exists := clients[clientIP]
		if !exists {
			limiter = &clientLimiter{
				limiter:    rate.NewLimiter(cfg.RequestsPerSecond, cfg.BurstSize),
				lastAccess: time.Now(),
			}
			clients[clientIP] = limiter
		} else {
			limiter.lastAccess = time.Now()
		}

		mu.Unlock()

		// check if request is allowed by rate limiter
		if !limiter.limiter.Allow() {
			utils.HandleErrorWithStatus(c, logger, http.StatusTooManyRequests, apierrors.ErrTooManyRequests)
			c.Abort()
			return
		}
		c.Next()
	}

}

// cleanupClients periodically cleans up old clients that haven't been used recently
func cleanupClients(clients map[string]*clientLimiter, mu *sync.Mutex, cleanupInterval time.Duration) {
	for {
		time.Sleep(time.Minute)

		mu.Lock()
		for ip, limiter := range clients {
			if time.Since(limiter.lastAccess) > cleanupInterval {
				delete(clients, ip)
			}
		}
		mu.Unlock()
	}
}
