package storage

type Config struct {
	URI           string `env:"MONGO_CRAWLER_URI"`
	User          string `env:"MONGO_USER"`
	Password      string `env:"MONGO_PASSWORD"`
	MaxPoolSize   uint64 `env:"MAX_POOL_SIZE" envDefault:"8"`
	MaxConnecting uint64 `env:"MAX_CONNECTING" envDefault:"8"`
}
