package metrics

import (
	"context"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

const (
	namespace = "spacebox_crawler"
)

// startScraping fills metrics in prometheus in background
func (m *Metrics) startScraping() {
	m.log.Info().Msg("start metrics scraper")

	var (
		// count of blocks for each status
		statusMetric = promauto.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "blocks_by_status",
			Help:      "Total blocks for each status",
		}, []string{"status"})

		// last processed height
		heightMetric = promauto.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "last_processed_block_height",
			Help:      "Last processed block height",
		})

		statusMap map[string]int
		ctx       = context.Background()
		ticker    = time.NewTicker(1 * time.Minute)
	)

	for {
		select {
		case <-m.stopScraping:
			m.log.Info().Msg("stop metrics scraper")
			return
		case <-ticker.C:
			if err := m.storage.Ping(ctx); err != nil {
				continue
			}
			blocks, err := m.storage.GetAllBlocks(ctx)
			if err != nil {
				m.log.Error().Err(err).Msg("can't get all blocks from storage")
				continue
			}

			statusMap = make(map[string]int) // clear map
			var maxHeight int64
			for _, b := range blocks {
				if b.Status.IsProcessed() && b.Height > maxHeight {
					maxHeight = b.Height
				}
				statusMap[b.Status.ToString()] += 1
			}

			heightMetric.Set(float64(maxHeight))

			for statusName, count := range statusMap {
				statusMetric.With(prometheus.Labels{"status": statusName}).Set(float64(count))
			}
		}
	}
}
