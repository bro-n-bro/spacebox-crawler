package cache

type Option[K, V comparable] func(*Cache[K, V])

func WithMetrics[K, V comparable](instanceName string) Option[K, V] {
	return func(c *Cache[K, V]) {
		c.specialMetric = instanceHitMissMetric
		c.globalMetric = hitMissMetric
		c.lengthMetric = cacheLenghtMetric
		c.instance = instanceName
	}
}

func WithCompareFunc[K, V comparable](compareFn CompareFn[V]) Option[K, V] {
	return func(c *Cache[K, V]) {
		c.compareFn = compareFn
	}
}
