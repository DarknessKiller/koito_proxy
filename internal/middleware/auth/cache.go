package auth

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
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

func (c *Cache) hashKey(key string) string {
	h := sha256.Sum256([]byte(key))
	return hex.EncodeToString(h[:])
}

func (c *Cache) Get(key string) (bool, bool) {
	hk := c.hashKey(key)
	c.mu.RLock()
	item, ok := c.items[hk]
	c.mu.RUnlock()

	if !ok {
		return false, false
	}

	if !time.Now().After(item.ExpiresAt) {
		return true, true
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	item, ok = c.items[hk]
	if ok && time.Now().After(item.ExpiresAt) {
		delete(c.items, hk)
	}

	return false, false
}

func (c *Cache) Set(key string, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.items[c.hashKey(key)] = CacheItem{
		ExpiresAt: time.Now().Add(ttl),
	}
}

func (c *Cache) StartCleanup(ctx context.Context, interval time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				c.evictExpired()
			}
		}
	}()
}

func (c *Cache) evictExpired() {
	c.mu.Lock()
	defer c.mu.Unlock()

	now := time.Now()
	for k, v := range c.items {
		if now.After(v.ExpiresAt) {
			delete(c.items, k)
		}
	}
}
