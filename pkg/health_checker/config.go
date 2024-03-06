package healthchecker

import "time"

type Config struct {
	MaxBlockLag  time.Duration `env:"HEALTHCHECK_MAX_LAST_BLOCK_LAG" envDefault:"5m"`
	Interval     time.Duration `env:"HEALTHCHECK_INTERVAL" envDefault:"1m"`
	StartDelay   time.Duration `env:"HEALTHCHECK_START_DELAY"`
	Enabled      bool          `env:"HEALTHCHECK_ENABLED" envDefault:"false"`
	FatalOnCheck bool          `env:"HEALTHCHECK_FATAL_ON_CHECK" envDefault:"true"`
}
