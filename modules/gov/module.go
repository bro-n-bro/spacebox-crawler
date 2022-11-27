package bank

import (
	"bro-n-bro-osmosis/internal/rep"
	tb "bro-n-bro-osmosis/pkg/mapper/to_broker"
	"bro-n-bro-osmosis/types"
	"os"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/rs/zerolog"

	grpcClient "bro-n-bro-osmosis/client/grpc"
)

var (
	_ types.Module        = &Module{}
	_ types.GenesisModule = &Module{}
	_ types.BlockModule   = &Module{}
	_ types.MessageModule = &Module{}
)

type Module struct {
	log    *zerolog.Logger
	client *grpcClient.Client
	broker rep.Broker
	tbM    tb.ToBroker
	cdc    codec.Codec
}

func New(b rep.Broker, cli *grpcClient.Client, tbM tb.ToBroker, cdc codec.Codec) *Module {
	l := zerolog.New(os.Stderr).Output(zerolog.ConsoleWriter{Out: os.Stderr}).With().Timestamp().
		Str("module", "gov").Logger()

	return &Module{
		log:    &l,
		broker: b,
		client: cli,
		tbM:    tbM,
		cdc:    cdc,
	}
}

func (m *Module) Name() string { return "gov" }
