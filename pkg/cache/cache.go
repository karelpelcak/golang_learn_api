package cache

import (
	"sync"
	"time"
)

// Cache represents a simple in-memory cache
type Cache struct {
	items map[string]*cacheItem
	mutex sync.RWMutex
}

// cacheItem represents a cached item with expiration
type cacheItem struct {
	value      interface{}
	expiration time.Time
}

// NewCache creates a new cache instance
func NewCache() *Cache {
	cache := &Cache{
		items: make(map[string]*cacheItem),
	}
	
	// Start cleanup goroutine
	go cache.cleanup()
	
	return cache
}

// Set adds an item to the cache with a TTL (in seconds)
func (c *Cache) Set(key string, value interface{}, ttl int) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	
	c.items[key] = &cacheItem{
		value:      value,
		expiration: time.Now().Add(time.Duration(ttl) * time.Second),
	}
}

// Get retrieves an item from the cache
func (c *Cache) Get(key string) (interface{}, bool) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	
	item, exists := c.items[key]
	if !exists {
		return nil, false
	}
	
	// Check if item has expired
	if time.Now().After(item.expiration) {
		// Item expired, remove it
		go c.delete(key)
		return nil, false
	}
	
	return item.value, true
}

// Delete removes an item from the cache
func (c *Cache) Delete(key string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	
	delete(c.items, key)
}

// delete removes an item from the cache (internal method without locking)
func (c *Cache) delete(key string) {
	delete(c.items, key)
}

// cleanup periodically removes expired items
func (c *Cache) cleanup() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()
	
	for range ticker.C {
		c.mutex.Lock()
		now := time.Now()
		for key, item := range c.items {
			if now.After(item.expiration) {
				delete(c.items, key)
			}
		}
		c.mutex.Unlock()
	}
}