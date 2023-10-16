package slashing

import (
	"github.com/rs/zerolog"

	grpcClient "github.com/bro-n-bro/spacebox-crawler/client/grpc"
	"github.com/bro-n-bro/spacebox-crawler/modules/utils"
	tb "github.com/bro-n-bro/spacebox-crawler/pkg/mapper/to_broker"
	"github.com/bro-n-bro/spacebox-crawler/types"
)

const (
	ModuleName = "slashing"
)

var (
	_ types.Module              = &Module{}
	_ types.BlockHandler        = &Module{}
	_ types.MessageHandler      = &Module{}
	_ types.BeginBlockerHandler = &Module{}
)

type Module struct {
	log    *zerolog.Logger
	client *grpcClient.Client
	tbM    tb.ToBroker
	broker broker
}

func New(b broker, cli *grpcClient.Client, tbM tb.ToBroker) *Module {
	return &Module{
		log:    utils.NewModuleLogger(ModuleName),
		client: cli,
		tbM:    tbM,
		broker: b,
	}
}

func (m *Module) Name() string { return ModuleName }
