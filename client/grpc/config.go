package grpc

import "time"

type (
	Config struct {
		Host             string        `env:"GRPC_URL" envDefault:"http://localhost:9090"`
		SecureConnection bool          `env:"GRPC_SECURE_CONNECTION" envDefault:"false"`
		MetricsEnabled   bool          `env:"METRICS_ENABLED" envDefault:"false"`
		Timeout          time.Duration `env:"GRPC_TIMEOUT" envDefault:"15s"`
	}
)
