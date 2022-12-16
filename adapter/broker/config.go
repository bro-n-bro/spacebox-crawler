package broker

type Config struct {
	ServerURL string `env:"BROKER_SERVER"`
	Enabled   bool   `env:"BROKER_ENABLED"`
}
