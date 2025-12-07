package handlers

import (
	"net/http"
	"runtime"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/maxqr-api/internal/cache"
)

// HealthHandler handles health check requests
type HealthHandler struct {
	startTime time.Time
	cache     *cache.Cache
}

// NewHealthHandler creates a new health handler
func NewHealthHandler(qrCache *cache.Cache) *HealthHandler {
	return &HealthHandler{
		startTime: time.Now(),
		cache:     qrCache,
	}
}

// Health handles GET /health
func (h *HealthHandler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":    "healthy",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

// Ready handles GET /ready
func (h *HealthHandler) Ready(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "ready",
	})
}

// Stats handles GET /stats
func (h *HealthHandler) Stats(c *gin.Context) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	uptime := time.Since(h.startTime)

	cacheStats := cache.Stats{}
	if h.cache != nil {
		cacheStats = h.cache.Stats()
	}

	c.JSON(http.StatusOK, gin.H{
		"uptime_seconds": uptime.Seconds(),
		"uptime_human":   uptime.String(),
		"memory": gin.H{
			"alloc_mb":       m.Alloc / 1024 / 1024,
			"total_alloc_mb": m.TotalAlloc / 1024 / 1024,
			"sys_mb":         m.Sys / 1024 / 1024,
			"num_gc":         m.NumGC,
		},
		"goroutines": runtime.NumGoroutine(),
		"cache":      cacheStats,
	})
}
