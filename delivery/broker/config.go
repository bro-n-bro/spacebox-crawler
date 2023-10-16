package broker

type (
	Config struct {
		ServerURL       string `env:"BROKER_SERVER"`
		PartitionsCount int    `env:"PARTITIONS_COUNT" envDefault:"1"`
		Enabled         bool   `env:"BROKER_ENABLED"`
	}
)
