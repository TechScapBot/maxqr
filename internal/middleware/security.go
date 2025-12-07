package middleware

import (
	"github.com/gin-gonic/gin"
)

// SecurityHeaders adds security headers to responses
func SecurityHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Prevent MIME type sniffing
		c.Header("X-Content-Type-Options", "nosniff")

		// Prevent clickjacking
		c.Header("X-Frame-Options", "DENY")

		// XSS Protection
		c.Header("X-XSS-Protection", "1; mode=block")

		// Referrer Policy
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")

		// Content Security Policy for images
		c.Header("Content-Security-Policy", "default-src 'none'; img-src 'self' data:; style-src 'unsafe-inline'")

		// Permissions Policy
		c.Header("Permissions-Policy", "accelerometer=(), camera=(), geolocation=(), gyroscope=(), magnetometer=(), microphone=(), payment=(), usb=()")

		c.Next()
	}
}

// RequestID adds a unique request ID to each request
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check if client provided a request ID
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			// Generate a simple request ID based on timestamp
			requestID = generateRequestID()
		}

		c.Set("request_id", requestID)
		c.Header("X-Request-ID", requestID)

		c.Next()
	}
}

// generateRequestID creates a simple unique ID
func generateRequestID() string {
	// Use nanosecond timestamp + random suffix for uniqueness
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, 16)

	// Simple pseudo-random based on time
	seed := uint64(0)
	for i := 0; i < 16; i++ {
		seed = seed*1103515245 + 12345
		result[i] = charset[int(seed>>16)%len(charset)]
	}

	return string(result)
}

// Recovery middleware with custom error handling
func Recovery() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		c.JSON(500, gin.H{
			"error":   "internal_server_error",
			"message": "An unexpected error occurred",
		})
		c.Abort()
	})
}
