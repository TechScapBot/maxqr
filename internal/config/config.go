package config

import (
	"os"
	"strconv"
	"time"
)

// Config holds all application configuration
type Config struct {
	// Server
	Port            string
	Host            string
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	ShutdownTimeout time.Duration

	// Security
	APIKey          string
	AllowedOrigins  string
	EnableCORS      bool

	// Rate Limiting
	RateLimitEnabled    bool
	RateLimitPerSecond  int
	RateLimitBurst      int

	// Cache
	CacheEnabled        bool
	CacheMaxSizeMB      int
	CacheDefaultTTL     time.Duration
	CacheCleanupInterval time.Duration

	// QR Generation
	QRDefaultSize       int
	QRMaxSize           int
	QRDefaultRecovery   string

	// Logging
	LogLevel        string
	LogFormat       string
}

// Load loads configuration from environment variables
func Load() *Config {
	return &Config{
		// Server
		Port:            getEnv("PORT", "8080"),
		Host:            getEnv("HOST", "0.0.0.0"),
		ReadTimeout:     getDuration("READ_TIMEOUT", 10*time.Second),
		WriteTimeout:    getDuration("WRITE_TIMEOUT", 30*time.Second),
		ShutdownTimeout: getDuration("SHUTDOWN_TIMEOUT", 30*time.Second),

		// Security
		APIKey:         getEnv("API_KEY", ""),
		AllowedOrigins: getEnv("ALLOWED_ORIGINS", "*"),
		EnableCORS:     getBool("ENABLE_CORS", true),

		// Rate Limiting
		RateLimitEnabled:   getBool("RATE_LIMIT_ENABLED", true),
		RateLimitPerSecond: getInt("RATE_LIMIT_PER_SECOND", 100),
		RateLimitBurst:     getInt("RATE_LIMIT_BURST", 200),

		// Cache
		CacheEnabled:         getBool("CACHE_ENABLED", true),
		CacheMaxSizeMB:       getInt("CACHE_MAX_SIZE_MB", 100),
		CacheDefaultTTL:      getDuration("CACHE_DEFAULT_TTL", 5*time.Minute),
		CacheCleanupInterval: getDuration("CACHE_CLEANUP_INTERVAL", 10*time.Minute),

		// QR Generation
		QRDefaultSize:     getInt("QR_DEFAULT_SIZE", 300),
		QRMaxSize:         getInt("QR_MAX_SIZE", 1000),
		QRDefaultRecovery: getEnv("QR_DEFAULT_RECOVERY", "M"),

		// Logging
		LogLevel:  getEnv("LOG_LEVEL", "info"),
		LogFormat: getEnv("LOG_FORMAT", "json"),
	}
}

// Helper functions
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return defaultValue
}

func getBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolVal, err := strconv.ParseBool(value); err == nil {
			return boolVal
		}
	}
	return defaultValue
}

func getDuration(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}
