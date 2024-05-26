package cache

import (
	"github.com/patrickmn/go-cache"
	"time"
)

// CacheInMem is a cache implementation.
type CacheInMem struct {
	cache *cache.Cache
}

// GetDel pools a value from the cache.
func (p *CacheInMem) GetDel(key string) (any, bool) {
	value, ok := p.Get(key)
	if !ok {
		return nil, false
	}
	p.cache.Delete(key)
	return value, true
}

// Get gets a value from the cache.
func (p *CacheInMem) Get(key string) (value any, ok bool) {
	return p.cache.Get(key)
}

// Put sets a value in the cache.
func (p *CacheInMem) Put(key string, value any, ttl time.Duration) {
	p.cache.Set(key, value, ttl)
}

// New creates a new cache instance.
func New() Cache {
	return &CacheInMem{
		cache: cache.New(cache.NoExpiration, cache.NoExpiration),
	}
}
