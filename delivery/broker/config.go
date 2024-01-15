package broker

type (
	Config struct {
		ServerURL       string `env:"BROKER_SERVER"`
		PartitionsCount int    `env:"PARTITIONS_COUNT" envDefault:"1"`
		MaxMessageBytes int    `env:"MAX_MESSAGE_MAX_BYTES" envDefault:"5242880"` // 5MB
		Enabled         bool   `env:"BROKER_ENABLED"`
	}
)
