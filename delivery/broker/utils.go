package broker

import (
	lru "github.com/hashicorp/golang-lru/v2"
)

func updateCacheValue[K, V comparable](cache *lru.Cache[K, V], key K, newVal V,
	compareFn func(curVal V) bool) (updated bool) {

	cacheVal, ok := cache.Get(key)
	if !ok {
		cache.Add(key, newVal)
		return true
	}

	if compareFn(cacheVal) {
		cache.Add(key, newVal)
		return true
	}

	return
}
