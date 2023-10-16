package feegrant

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/rs/zerolog"

	grpcClient "github.com/bro-n-bro/spacebox-crawler/client/grpc"
	"github.com/bro-n-bro/spacebox-crawler/modules/utils"
	tb "github.com/bro-n-bro/spacebox-crawler/pkg/mapper/to_broker"
	"github.com/bro-n-bro/spacebox-crawler/types"
)

const (
	ModuleName = "feegrant"
)

var (
	_ types.Module         = &Module{}
	_ types.MessageHandler = &Module{}
)

type Module struct {
	log    *zerolog.Logger
	client *grpcClient.Client
	broker broker
	tbM    tb.ToBroker
	cdc    codec.Codec
}

func New(b broker, cli *grpcClient.Client, tbM tb.ToBroker, cdc codec.Codec) *Module {
	return &Module{
		log:    utils.NewModuleLogger(ModuleName),
		broker: b,
		client: cli,
		tbM:    tbM,
		cdc:    cdc,
	}
}

func (m *Module) Name() string { return ModuleName }
