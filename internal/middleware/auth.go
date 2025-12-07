package middleware

import (
	"crypto/subtle"
	"net/http"

	"github.com/gin-gonic/gin"
)

// APIKeyAuth returns middleware for API key authentication
func APIKeyAuth(apiKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Skip auth if no API key is configured
		if apiKey == "" {
			c.Next()
			return
		}

		// Check X-API-Key header
		providedKey := c.GetHeader("X-API-Key")
		if providedKey == "" {
			// Also check query parameter as fallback
			providedKey = c.Query("api_key")
		}

		if providedKey == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "unauthorized",
				"message": "API key is required. Provide X-API-Key header or api_key query parameter.",
			})
			c.Abort()
			return
		}

		// Constant-time comparison to prevent timing attacks
		if subtle.ConstantTimeCompare([]byte(providedKey), []byte(apiKey)) != 1 {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "unauthorized",
				"message": "Invalid API key",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// OptionalAPIKeyAuth allows requests without API key but records if key is present
func OptionalAPIKeyAuth(apiKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		providedKey := c.GetHeader("X-API-Key")
		if providedKey == "" {
			providedKey = c.Query("api_key")
		}

		if providedKey != "" && apiKey != "" {
			if subtle.ConstantTimeCompare([]byte(providedKey), []byte(apiKey)) == 1 {
				c.Set("authenticated", true)
			} else {
				c.Set("authenticated", false)
			}
		} else {
			c.Set("authenticated", false)
		}

		c.Next()
	}
}
