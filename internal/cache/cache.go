package cache

import (
	"sync"
	"sync/atomic"
	"time"

	gocache "github.com/patrickmn/go-cache"
)

// Cache is a high-performance in-memory cache for QR codes
type Cache struct {
	store       *gocache.Cache
	hitCount    uint64 // Use atomic operations - no mutex needed
	missCount   uint64 // Use atomic operations - no mutex needed
	mu          sync.RWMutex
	maxSize     int
	currentSize int64 // Use int64 for atomic operations
}

// Config holds cache configuration
type Config struct {
	DefaultExpiration time.Duration
	CleanupInterval   time.Duration
	MaxSizeBytes      int // Maximum cache size in bytes (0 = unlimited)
}

// DefaultConfig returns default cache configuration
func DefaultConfig() Config {
	return Config{
		DefaultExpiration: 5 * time.Minute,
		CleanupInterval:   10 * time.Minute,
		MaxSizeBytes:      100 * 1024 * 1024, // 100MB default
	}
}

// NewCache creates a new cache instance
func NewCache(cfg Config) *Cache {
	return &Cache{
		store:   gocache.New(cfg.DefaultExpiration, cfg.CleanupInterval),
		maxSize: cfg.MaxSizeBytes,
	}
}

// Get retrieves an item from the cache (lock-free for counters)
func (c *Cache) Get(key string) ([]byte, bool) {
	if item, found := c.store.Get(key); found {
		atomic.AddUint64(&c.hitCount, 1) // Lock-free atomic increment
		return item.([]byte), true
	}
	atomic.AddUint64(&c.missCount, 1) // Lock-free atomic increment
	return nil, false
}

// Set stores an item in the cache
func (c *Cache) Set(key string, value []byte) {
	c.SetWithExpiration(key, value, gocache.DefaultExpiration)
}

// SetWithExpiration stores an item with a specific expiration
func (c *Cache) SetWithExpiration(key string, value []byte, expiration time.Duration) {
	size := int64(len(value))

	// Check if we need to evict items
	if c.maxSize > 0 {
		c.mu.Lock()
		for c.currentSize+size > int64(c.maxSize) && c.store.ItemCount() > 0 {
			// Simple eviction: delete oldest items
			c.evictOne()
		}
		c.currentSize += size
		c.mu.Unlock()
	}

	c.store.Set(key, value, expiration)
}

// evictOne removes one item from cache (must be called with lock held)
func (c *Cache) evictOne() {
	items := c.store.Items()
	var oldestKey string
	var oldestTime int64 = time.Now().UnixNano()

	for key, item := range items {
		if item.Expiration < oldestTime && item.Expiration > 0 {
			oldestTime = item.Expiration
			oldestKey = key
		}
	}

	if oldestKey != "" {
		if item, found := c.store.Get(oldestKey); found {
			c.currentSize -= int64(len(item.([]byte)))
		}
		c.store.Delete(oldestKey)
	}
}

// Delete removes an item from the cache
func (c *Cache) Delete(key string) {
	if item, found := c.store.Get(key); found {
		c.mu.Lock()
		c.currentSize -= int64(len(item.([]byte)))
		c.mu.Unlock()
	}
	c.store.Delete(key)
}

// Clear removes all items from the cache
func (c *Cache) Clear() {
	c.store.Flush()
	c.mu.Lock()
	c.currentSize = 0
	c.hitCount = 0
	c.missCount = 0
	c.mu.Unlock()
}

// Stats returns cache statistics
type Stats struct {
	ItemCount   int     `json:"item_count"`
	HitCount    uint64  `json:"hit_count"`
	MissCount   uint64  `json:"miss_count"`
	HitRate     float64 `json:"hit_rate"`
	SizeBytes   int64   `json:"size_bytes"`
	MaxBytes    int     `json:"max_bytes"`
}

// Stats returns current cache statistics (lock-free for counters)
func (c *Cache) Stats() Stats {
	// Use atomic loads for counters - no lock needed
	hitCount := atomic.LoadUint64(&c.hitCount)
	missCount := atomic.LoadUint64(&c.missCount)

	total := hitCount + missCount
	hitRate := float64(0)
	if total > 0 {
		hitRate = float64(hitCount) / float64(total)
	}

	c.mu.RLock()
	currentSize := c.currentSize
	c.mu.RUnlock()

	return Stats{
		ItemCount: c.store.ItemCount(),
		HitCount:  hitCount,
		MissCount: missCount,
		HitRate:   hitRate,
		SizeBytes: currentSize,
		MaxBytes:  c.maxSize,
	}
}

// ItemCount returns the number of items in the cache
func (c *Cache) ItemCount() int {
	return c.store.ItemCount()
}

// Global cache instance
var globalCache *Cache
var once sync.Once

// Global returns the global cache instance
func Global() *Cache {
	once.Do(func() {
		globalCache = NewCache(DefaultConfig())
	})
	return globalCache
}

// InitGlobal initializes the global cache with custom config
func InitGlobal(cfg Config) {
	once.Do(func() {
		globalCache = NewCache(cfg)
	})
}
