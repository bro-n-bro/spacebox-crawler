package bank

import (
	"bro-n-bro-osmosis/types"
	"os"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/rs/zerolog"

	"bro-n-bro-osmosis/adapter/broker"
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
	broker *broker.Broker
	client *grpcClient.Client
	cdc    codec.Codec
}

func New(b *broker.Broker, cli *grpcClient.Client, cdc codec.Codec) *Module {
	l := zerolog.New(os.Stderr).Output(zerolog.ConsoleWriter{Out: os.Stderr}).With().Timestamp().
		Str("module", "gov").Logger()

	return &Module{
		log:    &l,
		broker: b,
		client: cli,
		cdc:    cdc,
	}
}

func (m *Module) Name() string { return "gov" }
