package server

import (
	"context"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"

	"github.com/bro-n-bro/spacebox-crawler/adapter/storage/model"
)

const (
	namespace = "spacebox_crawler"
)

// startMetricsScrapper fills metrics in prometheus in background.
func (s *Server) startMetricsScrapper() {
	s.log.Info().Msg("start metrics scraper")

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

		// count of error messages in storage
		errorMessagesCount = promauto.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "total_error_messages",
			Help:      "Total error messages",
		})

		statusMap map[string]int
		ctx       = context.Background()
		ticker    = time.NewTicker(1 * time.Minute)
	)

	defer ticker.Stop()

	for {
		select {
		case <-s.stopScraping:
			s.log.Info().Msg("stop metrics scraper")
			return
		case <-ticker.C:
			if err := s.storage.Ping(ctx); err != nil {
				continue
			}
			blocks, err := s.storage.GetAllBlocks(ctx)
			if err != nil {
				s.log.Error().Err(err).Msg("can't get all blocks from storage")
				continue
			}

			statusMap = map[string]int{
				model.StatusProcessing.ToString(): 0,
				model.StatusProcessed.ToString():  0,
				model.StatusError.ToString():      0,
			}

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

			count, err := s.storage.CountErrorMessage(ctx)
			if err != nil {
				s.log.Error().Err(err).Msg("can't get count of error messages")
				continue
			}
			errorMessagesCount.Set(float64(count))
		}
	}
}
