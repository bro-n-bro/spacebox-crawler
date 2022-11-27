package auth

import (
	"bro-n-bro-osmosis/internal/rep"
	"os"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/rs/zerolog"

	grpcClient "bro-n-bro-osmosis/client/grpc"
	"bro-n-bro-osmosis/modules/messages"
	tb "bro-n-bro-osmosis/pkg/mapper/to_broker"
	"bro-n-bro-osmosis/types"
)

var (
	_ types.Module        = &Module{}
	_ types.GenesisModule = &Module{}
	_ types.MessageModule = &Module{}
)

type Module struct {
	log    *zerolog.Logger
	broker rep.Broker
	client *grpcClient.Client
	tbM    tb.ToBroker
	cdc    codec.Codec
	parser messages.MessageAddressesParser
}

func New(b rep.Broker, cli *grpcClient.Client, tb tb.ToBroker,
	cdc codec.Codec, parser messages.MessageAddressesParser) *Module {

	l := zerolog.New(os.Stderr).Output(zerolog.ConsoleWriter{Out: os.Stderr}).With().Timestamp().
		Str("module", "auth").Logger()

	return &Module{
		log:    &l,
		broker: b,
		tbM:    tb,
		client: cli,
		cdc:    cdc,
		parser: parser,
	}
}

func (m *Module) Name() string { return "auth" }
