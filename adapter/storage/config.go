package storage

type Config struct {
	URI      string `env:"MONGO_URI"`
	User     string `env:"MONGO_USER"`
	Password string `env:"MONGO_PASSWORD"`
}
