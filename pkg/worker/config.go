package worker

import "time"

type Config struct {
	ProcessErrorBlocksInterval time.Duration `env:"PROCESS_ERROR_BLOCKS_INTERVAL" envDefault:"1m"`
	ProcessNewBlocks           bool          `env:"SUBSCRIBE_NEW_BLOCKS"` // FIXME: or use ws enabled???
	ProcessErrorBlocks         bool          `env:"PROCESS_ERROR_BLOCKS" envDefault:"true"`
	MetricsEnabled             bool          `env:"WORKER_METRICS_ENABLED" envDefault:"false"`
	RecoveryMode               bool          `env:"RECOVERY_MODE"`
	WorkersCount               int           `env:"WORKERS_COUNT" envDefault:"1"`
	StartHeight                int64         `env:"START_HEIGHT" envDefault:"-1"`
	StopHeight                 int64         `env:"STOP_HEIGHT"`
}
