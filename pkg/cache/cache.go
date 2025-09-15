package cache

import (
	"sync"
	"time"
)

type Cache struct {
	items map[string]*cacheItem
	mutex sync.RWMutex
}

type cacheItem struct {
	value      interface{}
	expiration time.Time
}

func NewCache() *Cache {
	cache := &Cache{
		items: make(map[string]*cacheItem),
	}

	go cache.cleanup()

	return cache
}

func (c *Cache) Set(key string, value interface{}, ttl int) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.items[key] = &cacheItem{
		value:      value,
		expiration: time.Now().Add(time.Duration(ttl) * time.Second),
	}
}

func (c *Cache) Get(key string) (interface{}, bool) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	item, exists := c.items[key]
	if !exists {
		return nil, false
	}

	if time.Now().After(item.expiration) {
		go c.delete(key)
		return nil, false
	}

	return item.value, true
}

func (c *Cache) Delete(key string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	delete(c.items, key)
}

func (c *Cache) delete(key string) {
	delete(c.items, key)
}

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
