package worker

import (
	"time"
)

func (w *Worker) withMetrics(labelType string, fn func() error) error {
	if !w.cfg.MetricsEnabled || w.metrics == nil {
		return fn()
	}

	start := time.Now()
	err := fn()
	w.metrics.durMetric.WithLabelValues(labelType).Observe(time.Since(start).Seconds())

	return err
}
