package grid

import (
	"github.com/rs/zerolog"

	grpcClient "github.com/bro-n-bro/spacebox-crawler/client/grpc"
	"github.com/bro-n-bro/spacebox-crawler/modules/utils"
	tb "github.com/bro-n-bro/spacebox-crawler/pkg/mapper/to_broker"
	"github.com/bro-n-bro/spacebox-crawler/types"
)

const (
	ModuleName = "grid"
)

var (
	_ types.Module       = &Module{}
	_ types.BlockHandler = &Module{}
)

type Module struct {
	log    *zerolog.Logger
	client *grpcClient.Client
	broker broker
	tbM    tb.ToBroker
}

func New(b broker, cli *grpcClient.Client, tbM tb.ToBroker) *Module {
	return &Module{
		log:    utils.NewModuleLogger(ModuleName),
		broker: b,
		client: cli,
		tbM:    tbM,
	}
}

func (m *Module) Name() string { return ModuleName }
