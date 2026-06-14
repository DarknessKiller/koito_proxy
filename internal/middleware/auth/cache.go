package auth

import (
	"sync"
	"time"
)

type CacheItem struct {
	ExpiresAt time.Time
}

type Cache struct {
	mu    sync.RWMutex
	items map[string]CacheItem
}

func NewCache() *Cache {
	return &Cache{
		items: make(map[string]CacheItem),
	}
}

func (c *Cache) Get(key string) (bool, bool) {
	c.mu.RLock()
	item, ok := c.items[key]
	c.mu.RUnlock()

	if !ok {
		return false, false
	}

	if !time.Now().After(item.ExpiresAt) {
		return true, true
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	item, ok = c.items[key]
	if ok && time.Now().After(item.ExpiresAt) {
		delete(c.items, key)
	}

	return false, false
}

func (c *Cache) Set(key string, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.items[key] = CacheItem{
		ExpiresAt: time.Now().Add(ttl),
	}
}
