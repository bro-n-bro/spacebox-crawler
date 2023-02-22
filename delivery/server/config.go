package server

type Config struct {
	Port           string `env:"SERVER_PORT" envDefault:"8080"`
	MetricsEnabled bool   `env:"METRICS_ENABLED" envDefault:"false"`
}
