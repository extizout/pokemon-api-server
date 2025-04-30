package cache

import (
	"time"

	"github.com/patrickmn/go-cache"
)

type memoryCache struct {
	cache *cache.Cache
}

func NewMemoryCache(defaultExpiration, cleanupInterval time.Duration) CacheService {
	return &memoryCache{
		cache: cache.New(defaultExpiration, cleanupInterval),
	}
}

func (m *memoryCache) Get(key string) (interface{}, bool) {
	return m.cache.Get(key)
}

func (m *memoryCache) Set(key string, value interface{}, ttl time.Duration) {
	m.cache.Set(key, value, ttl)
}

func (m *memoryCache) SetWithDefaultExpiration(key string, value interface{}) {
	m.cache.Set(key, value, cache.DefaultExpiration)
}
