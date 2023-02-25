package broker

// checkCache checks the cache to see if it needs to do something with the new value.
func checkCache[K, V comparable](key K, val V, cache cacheI[K, V]) bool {
	if cache != nil && cache.UpdateCacheValue(key, val) {
		// value updated
		return true
	}

	return false
}
