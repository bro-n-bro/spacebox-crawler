package rpc

import "time"

type Config struct {
	Host           string        `env:"RPC_URL" envDefault:"http://localhost:26657"`
	MetricsEnabled bool          `env:"METRICS_ENABLED" envDefault:"false"`
	Timeout        time.Duration `env:"RPC_TIMEOUT" envDefault:"15s"`
}
