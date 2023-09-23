package cache

import (
	lru "github.com/hashicorp/golang-lru/v2"
	"github.com/prometheus/client_golang/prometheus"
)

type (
	CompareFn[V comparable] func(cacheVal, newVal V) bool

	Cache[K, V comparable] struct {
		cache *lru.Cache[K, V]

		globalMetric  *prometheus.CounterVec
		lengthMetric  *prometheus.GaugeVec
		specialMetric *prometheus.CounterVec

		compareFn CompareFn[V]
		instance  string
	}
)

// New creates the new Cache instance
func New[K, V comparable](size int, opts ...Option[K, V]) (*Cache[K, V], error) {
	c, err := lru.New[K, V](size)
	if err != nil {
		return nil, err
	}

	cache := &Cache[K, V]{cache: c}

	for _, opt := range opts {
		opt(cache)
	}

	return cache, nil
}

func (c *Cache[K, V]) Patch(opt Option[K, V]) {
	opt(c)
}

// UpdateCacheValue returns true if nev value updated in cache.
func (c *Cache[K, V]) UpdateCacheValue(key K, newVal V) (updated bool) {
	cacheVal, ok := c.cache.Get(key)
	if !ok {
		c.cache.Add(key, newVal)
		c.miss()
		return true
	}

	if c.compareFn != nil && c.compareFn(cacheVal, newVal) {
		c.cache.Add(key, newVal)
		c.miss()
		return true
	}

	c.hit()
	return
}

func (c *Cache[K, V]) hit() {
	if c.globalMetric != nil {
		c.globalMetric.WithLabelValues("hit").Inc()
	}

	if c.specialMetric != nil {
		c.specialMetric.With(prometheus.Labels{"instance": c.instance, "action": "hit"}).Inc()
	}
}

func (c *Cache[K, V]) miss() {
	if c.globalMetric != nil {
		c.globalMetric.WithLabelValues("miss").Inc()
	}
	if c.specialMetric != nil {
		c.specialMetric.With(prometheus.Labels{"instance": c.instance, "action": "miss"}).Inc()
	}
	if c.lengthMetric != nil {
		c.lengthMetric.With(prometheus.Labels{"instance": c.instance}).Set(float64(c.cache.Len()))
	}
}
