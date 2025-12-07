package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/maxqr-api/internal/cache"
	"github.com/maxqr-api/internal/config"
	"github.com/maxqr-api/internal/handlers"
	"github.com/maxqr-api/internal/middleware"
	"github.com/maxqr-api/internal/qrgen"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Set Gin mode
	if cfg.LogLevel == "debug" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	// Initialize cache
	qrCache := cache.NewCache(cache.Config{
		DefaultExpiration: cfg.CacheDefaultTTL,
		CleanupInterval:   cfg.CacheCleanupInterval,
		MaxSizeBytes:      cfg.CacheMaxSizeMB * 1024 * 1024,
	})

	// Initialize QR generator
	generator := qrgen.NewGenerator(qrgen.DefaultConfig())

	// Initialize handlers
	qrHandler := handlers.NewQRHandler(generator, qrCache, cfg.CacheEnabled)
	bankHandler := handlers.NewBankHandler()
	healthHandler := handlers.NewHealthHandler(qrCache)

	// Create Gin router
	router := gin.New()

	// Add middleware
	router.Use(gin.Logger())
	router.Use(middleware.Recovery())
	router.Use(middleware.SecurityHeaders())
	router.Use(middleware.RequestID())

	// CORS configuration
	if cfg.EnableCORS {
		corsConfig := cors.Config{
			AllowOrigins:     []string{"*"},
			AllowMethods:     []string{"GET", "POST", "OPTIONS"},
			AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "X-API-Key", "X-Request-ID"},
			ExposeHeaders:    []string{"Content-Length", "X-Request-ID"},
			AllowCredentials: false,
			MaxAge:           12 * time.Hour,
		}

		if cfg.AllowedOrigins != "*" {
			corsConfig.AllowOrigins = []string{cfg.AllowedOrigins}
			corsConfig.AllowCredentials = true
		}

		router.Use(cors.New(corsConfig))
	}

	// Rate limiting
	if cfg.RateLimitEnabled {
		limiter := middleware.NewRateLimiter(cfg.RateLimitPerSecond, cfg.RateLimitBurst)
		router.Use(limiter.Middleware())
	}

	// Health endpoints (no auth)
	router.GET("/health", healthHandler.Health)
	router.GET("/ready", healthHandler.Ready)
	router.GET("/stats", healthHandler.Stats)

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Optional API key auth for enhanced limits
		if cfg.APIKey != "" {
			v1.Use(middleware.OptionalAPIKeyAuth(cfg.APIKey))
		}

		// Bank endpoints
		v1.GET("/banks", bankHandler.ListBanks)
		v1.GET("/banks/search", bankHandler.SearchBanks)
		v1.GET("/banks/:identifier", bankHandler.GetBank)

		// QR generation endpoints
		v1.POST("/generate", qrHandler.Generate)
		v1.GET("/quick", qrHandler.QuickGenerate)
		v1.GET("/qr/:bank_bin/:account_number", qrHandler.GenerateImage)

		// QR decode endpoint
		v1.POST("/decode", qrHandler.Decode)
	}

	// Legacy/simple endpoints for compatibility
	router.GET("/qr", qrHandler.QuickGenerate)

	// Create HTTP server
	server := &http.Server{
		Addr:         cfg.Host + ":" + cfg.Port,
		Handler:      router,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
	}

	// Start server in goroutine
	go func() {
		log.Printf("🚀 MaxQR API server starting on %s:%s", cfg.Host, cfg.Port)
		log.Printf("📊 Cache enabled: %v, Max size: %dMB", cfg.CacheEnabled, cfg.CacheMaxSizeMB)
		log.Printf("🔒 Rate limiting enabled: %v, Limit: %d req/s", cfg.RateLimitEnabled, cfg.RateLimitPerSecond)

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited gracefully")
}
