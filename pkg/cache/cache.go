package cache

import "time"

type (
	CacheService interface {
		Get(key string) (interface{}, bool)
		Set(key string, value interface{}, ttl time.Duration)
		SetWithDefaultExpiration(key string, value interface{})
	}
)
