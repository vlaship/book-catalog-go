package cache

import "time"

// Cache is a cache interface.
//
//go:generate mockgen -destination=../../test/mock/cache/mock-cache.go -package=mock . Cache
type Cache interface {
	Get(key string) (any, bool)
	GetDel(key string) (any, bool)
	Put(key string, value any, ttl time.Duration)
}
