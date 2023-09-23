package cache

import (
	"sync"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	instanceHitMissMetric *prometheus.CounterVec
	cacheLenghtMetric     *prometheus.GaugeVec
	hitMissMetric         *prometheus.CounterVec

	once sync.Once
)

func RegisterMetrics(namespace string) {
	once.Do(func() {
		instanceHitMissMetric = promauto.NewCounterVec(prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "special_cache",
			Help:      "Hit miss cache metric by instance",
		}, []string{"instance", "action"})

		cacheLenghtMetric = promauto.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "cache_length",
			Help:      "Items in cache by instance",
		}, []string{"instance"})

		hitMissMetric = promauto.NewCounterVec(prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "global_cache",
			Help:      "Hit miss cache metrics",
		}, []string{"action"})
	})
}
