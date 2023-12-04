package graph

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/rs/zerolog"

	grpcClient "github.com/bro-n-bro/spacebox-crawler/client/grpc"
	"github.com/bro-n-bro/spacebox-crawler/modules/utils"
	tb "github.com/bro-n-bro/spacebox-crawler/pkg/mapper/to_broker"
	"github.com/bro-n-bro/spacebox-crawler/types"
)

const (
	ModuleName = "graph"
)

var (
	_ types.Module         = &Module{}
	_ types.MessageHandler = &Module{}
)

type (
	Module struct {
		log    *zerolog.Logger
		broker broker
		tbM    tb.ToBroker
		cdc    codec.Codec
		cli    *grpcClient.Client
	}
)

func New(b broker, tbM tb.ToBroker, cdc codec.Codec, cli *grpcClient.Client) *Module {
	return &Module{
		log:    utils.NewModuleLogger(ModuleName),
		broker: b,
		tbM:    tbM,
		cdc:    cdc,
		cli:    cli,
	}
}

func (m *Module) Name() string { return ModuleName }
