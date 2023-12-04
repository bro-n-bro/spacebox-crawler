package rpc

import "time"

type Config struct {
	Host      string        `env:"RPC_URL" envDefault:"http://localhost:26657"`
	WSEnabled bool          `env:"WS_ENABLED" envDefault:"true"`
	Timeout   time.Duration `env:"RPC_TIMEOUT" envDefault:"15s"`
}
