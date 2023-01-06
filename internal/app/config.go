package app

import (
	"time"

	"github.com/hexy-dev/spacebox-crawler/adapter/broker"
	"github.com/hexy-dev/spacebox-crawler/adapter/storage"
	"github.com/hexy-dev/spacebox-crawler/client/grpc"
	"github.com/hexy-dev/spacebox-crawler/client/rpc"
	"github.com/hexy-dev/spacebox-crawler/pkg/worker"
)

type Config struct {
	StorageConfig storage.Config
	ChainPrefix   string   `env:"CHAIN_PREFIX"`
	LogLevel      string   `env:"LOG_LEVEL" envDefault:"info"`
	Modules       []string `env:"MODULES" required:"true"`
	GRPCConfig    grpc.Config
	RPCConfig     rpc.Config
	BrokerConfig  broker.Config
	WorkerConfig  worker.Config
	StartTimeout  time.Duration `env:"START_TIMEOUT"`
	StopTimeout   time.Duration `env:"STOP_TIMEOUT"`
}
