package messages

import (
	"bro-n-bro-osmosis/adapter/broker"
	grpcClient "bro-n-bro-osmosis/client/grpc"
	"bro-n-bro-osmosis/types"
	"os"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/rs/zerolog"
)

var (
	_ types.Module        = &Module{}
	_ types.MessageModule = &Module{}
)

type Module struct {
	log    *zerolog.Logger
	broker *broker.Broker
	client *grpcClient.Client
	cdc    codec.Codec
	parser MessageAddressesParser
}

func New(b *broker.Broker, cli *grpcClient.Client, cdc codec.Codec, parser MessageAddressesParser) *Module {
	l := zerolog.New(os.Stderr).Output(zerolog.ConsoleWriter{Out: os.Stderr}).With().Timestamp().
		Str("module", "bank").Logger()

	return &Module{
		log:    &l,
		broker: b,
		client: cli,
		cdc:    cdc,
		parser: parser,
	}
}

func (m *Module) Name() string { return "messages" }
