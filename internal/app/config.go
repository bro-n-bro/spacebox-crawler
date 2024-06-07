package app

import (
	"time"

	"github.com/bro-n-bro/spacebox-crawler/v2/adapter/storage"
	"github.com/bro-n-bro/spacebox-crawler/v2/client/grpc"
	"github.com/bro-n-bro/spacebox-crawler/v2/client/rpc"
	"github.com/bro-n-bro/spacebox-crawler/v2/delivery/broker"
	"github.com/bro-n-bro/spacebox-crawler/v2/delivery/server"
	healthchecker "github.com/bro-n-bro/spacebox-crawler/v2/pkg/health_checker"
	"github.com/bro-n-bro/spacebox-crawler/v2/pkg/worker"
)

type Config struct {
	ChainPrefix       string `env:"CHAIN_PREFIX"`
	DefaultDenom      string `env:"DEFAULT_DENOM" envDefault:"uatom"`
	LogLevel          string `env:"LOG_LEVEL" envDefault:"info"`
	Server            server.Config
	GRPCConfig        grpc.Config
	RPCConfig         rpc.Config
	BrokerConfig      broker.Config
	StorageConfig     storage.Config
	WorkerConfig      worker.Config
	HealthcheckConfig healthchecker.Config
	StartTimeout      time.Duration `env:"START_TIMEOUT"`
	StopTimeout       time.Duration `env:"STOP_TIMEOUT"`
	MetricsEnabled    bool          `env:"METRICS_ENABLED" envDefault:"false"`
}
