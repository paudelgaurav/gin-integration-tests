package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/paudelgaurav/gin-boilerplate/pkg/framework"
	"github.com/stretchr/testify/assert"
	"golang.org/x/time/rate"
)

// TestRateLimiterMiddleware tests the rate-limiting middleware behavior
func TestRateLimiterMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Configuration for rate limiter: 1 request per second with a burst of 1
	cfg := LimiterConfig{
		RequestsPerSecond: rate.Limit(1),
		BurstSize:         1,
		CleanupInterval:   1 * time.Minute,
	}

	logger := framework.GetLogger()
	middlewareFunc := NewRateLimitMiddleware(cfg, logger)

	// Initialize a test Gin router
	router := gin.New()
	router.Use(middlewareFunc)

	// Simple test route
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	// Create a test request
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)

	// Send the first request (should be successful)
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.JSONEq(t, `{"message":"success"}`, w.Body.String())

	// Send a second request immediately (should fail due to rate limiting)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusTooManyRequests, w.Code)
	assert.JSONEq(t, `{"error":"too many requests"}`, w.Body.String())

	// Wait for 1 second to reset the rate limiter
	time.Sleep(1 * time.Second)

	// Send another request after waiting (should succeed)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.JSONEq(t, `{"message":"success"}`, w.Body.String())
}
