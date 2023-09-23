package app

import (
	"time"

	"github.com/bro-n-bro/spacebox-crawler/adapter/storage"
	"github.com/bro-n-bro/spacebox-crawler/client/grpc"
	"github.com/bro-n-bro/spacebox-crawler/client/rpc"
	"github.com/bro-n-bro/spacebox-crawler/delivery/broker"
	"github.com/bro-n-bro/spacebox-crawler/delivery/server"
	"github.com/bro-n-bro/spacebox-crawler/pkg/worker"
)

type Config struct {
	ChainPrefix    string `env:"CHAIN_PREFIX"`
	DefaultDenom   string `env:"DEFAULT_DENOM" envDefault:"uatom"`
	LogLevel       string `env:"LOG_LEVEL" envDefault:"info"`
	Server         server.Config
	Modules        []string `env:"MODULES" required:"true"`
	GRPCConfig     grpc.Config
	RPCConfig      rpc.Config
	BrokerConfig   broker.Config
	StorageConfig  storage.Config
	WorkerConfig   worker.Config
	StartTimeout   time.Duration `env:"START_TIMEOUT"`
	StopTimeout    time.Duration `env:"STOP_TIMEOUT"`
	MetricsEnabled bool          `env:"METRICS_ENABLED" envDefault:"false"`
}
