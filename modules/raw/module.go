package raw

import (
	"github.com/rs/zerolog"

	rpcClient "github.com/bro-n-bro/spacebox-crawler/client/rpc"
	"github.com/bro-n-bro/spacebox-crawler/modules/utils"
	tb "github.com/bro-n-bro/spacebox-crawler/pkg/mapper/to_broker"
	"github.com/bro-n-bro/spacebox-crawler/types"
)

const (
	ModuleName = "raw"
)

var (
	_ types.Module       = &Module{}
	_ types.BlockHandler = &Module{}
)

type Module struct {
	log       *zerolog.Logger
	rpcClient *rpcClient.Client
	broker    broker
	tbM       tb.ToBroker
}

func New(b broker, cli *rpcClient.Client, tbM tb.ToBroker) *Module {
	return &Module{
		log:       utils.NewModuleLogger(ModuleName),
		broker:    b,
		rpcClient: cli,
		tbM:       tbM,
	}
}

func (m *Module) Name() string { return ModuleName }
