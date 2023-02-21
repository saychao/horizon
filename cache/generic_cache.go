package cache

import (
	"time"

	"github.com/cheekybits/genny/generic"
	cache "github.com/patrickmn/go-cache"
)

//go:generate genny -in=$GOFILE -out=asset_cache_generated.go gen "ValueType=core.Asset"

type ValueType generic.Type

type ValueTypeCache struct {
	_cache *cache.Cache
}

func newValueTypeCache(defaultExpiration, cleanupInterval time.Duration) *ValueTypeCache {
	return &ValueTypeCache{
		_cache: cache.New(defaultExpiration, cleanupInterval),
	}
}

func (c *ValueTypeCache) Set(key string, value ValueType, d time.Duration) {
	c._cache.Set(key, value, d)
}

func (c *ValueTypeCache) Get(k string) (result ValueType, ok bool) {
	rawResult, ok := c._cache.Get(k)
	if !ok {
		return
	}

	if rawResult == nil {
		return
	}

	result, ok = rawResult.(ValueType)
	return
}
