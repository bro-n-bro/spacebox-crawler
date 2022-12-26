package worker

type Config struct {
	ProcessNewBlocks   bool  `env:"SUBSCRIBE_NEW_BLOCKS"` // FIXME: or use ws enabled???
	ProcessErrorBlocks bool  `env:"PROCESS_ERROR_BLOCKS" envDefault:"true"`
	ChanSize           int   `env:"WORKER_CHAN_SIZE" envDefault:"8"`
	WorkersCount       int   `env:"WORKERS_COUNT" envDefault:"8"`
	StartHeight        int64 `env:"START_HEIGHT" envDefault:"1"`
	StopHeight         int64 `env:"STOP_HEIGHT" envDefault:"0"`
}
