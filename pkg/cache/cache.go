package cache

import (
	lru "github.com/hashicorp/golang-lru/v2"
)

type (
	CompareFn[V comparable] func(cacheVal, newVal V) bool

	Cache[K, V comparable] struct {
		cache     *lru.Cache[K, V]
		compareFn CompareFn[V]
	}
)

// New creates the new Cache instance
func New[K, V comparable](size int) (*Cache[K, V], error) {
	c, err := lru.New[K, V](size)
	if err != nil {
		return nil, err
	}
	return &Cache[K, V]{cache: c}, nil
}

func (c *Cache[K, V]) SetCompareFn(compareFn CompareFn[V]) {
	c.compareFn = compareFn
}

// UpdateCacheValue returns true if nev value updated in cache.
func (c *Cache[K, V]) UpdateCacheValue(key K, newVal V) (updated bool) {
	cacheVal, ok := c.cache.Get(key)
	if !ok {
		c.cache.Add(key, newVal)
		return true
	}

	if c.compareFn != nil && c.compareFn(cacheVal, newVal) {
		c.cache.Add(key, newVal)
		return true
	}

	return
}
