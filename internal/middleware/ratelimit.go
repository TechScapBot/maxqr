package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// RateLimiter holds rate limiting configuration
type RateLimiter struct {
	limiters sync.Map
	rate     rate.Limit
	burst    int
	cleanup  time.Duration
}

// limiterEntry represents a single client's rate limiter
type limiterEntry struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(perSecond int, burst int) *RateLimiter {
	rl := &RateLimiter{
		rate:    rate.Limit(perSecond),
		burst:   burst,
		cleanup: 5 * time.Minute,
	}

	// Start cleanup goroutine
	go rl.cleanupLoop()

	return rl
}

// getLimiter returns the rate limiter for a given key
func (rl *RateLimiter) getLimiter(key string) *rate.Limiter {
	if entry, exists := rl.limiters.Load(key); exists {
		e := entry.(*limiterEntry)
		e.lastSeen = time.Now()
		return e.limiter
	}

	limiter := rate.NewLimiter(rl.rate, rl.burst)
	entry := &limiterEntry{
		limiter:  limiter,
		lastSeen: time.Now(),
	}
	rl.limiters.Store(key, entry)
	return limiter
}

// cleanupLoop removes old limiters
func (rl *RateLimiter) cleanupLoop() {
	ticker := time.NewTicker(rl.cleanup)
	defer ticker.Stop()

	for range ticker.C {
		rl.cleanupOld()
	}
}

// cleanupOld removes limiters not seen recently
func (rl *RateLimiter) cleanupOld() {
	cutoff := time.Now().Add(-rl.cleanup)
	rl.limiters.Range(func(key, value interface{}) bool {
		entry := value.(*limiterEntry)
		if entry.lastSeen.Before(cutoff) {
			rl.limiters.Delete(key)
		}
		return true
	})
}

// Middleware returns a Gin middleware for rate limiting
func (rl *RateLimiter) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Use client IP as key
		key := c.ClientIP()

		limiter := rl.getLimiter(key)
		if !limiter.Allow() {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error":   "rate_limit_exceeded",
				"message": "Too many requests. Please try again later.",
				"retry_after": 1,
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// RateLimitByKey returns middleware that rate limits by a custom key extractor
func (rl *RateLimiter) RateLimitByKey(keyFunc func(*gin.Context) string) gin.HandlerFunc {
	return func(c *gin.Context) {
		key := keyFunc(c)

		limiter := rl.getLimiter(key)
		if !limiter.Allow() {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error":   "rate_limit_exceeded",
				"message": "Too many requests. Please try again later.",
				"retry_after": 1,
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// Global rate limiter instance
var globalLimiter *RateLimiter
var limiterOnce sync.Once

// GlobalLimiter returns the global rate limiter
func GlobalLimiter(perSecond, burst int) *RateLimiter {
	limiterOnce.Do(func() {
		globalLimiter = NewRateLimiter(perSecond, burst)
	})
	return globalLimiter
}
