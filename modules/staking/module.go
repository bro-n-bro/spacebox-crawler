package staking

import (
	"bro-n-bro-osmosis/types"
	"os"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/rs/zerolog"

	"bro-n-bro-osmosis/adapter/broker"
	grpcClient "bro-n-bro-osmosis/client/grpc"
	"bro-n-bro-osmosis/modules/messages"
)

var (
	_ types.Module        = &Module{}
	_ types.GenesisModule = &Module{}
	_ types.BlockModule   = &Module{}
	_ types.MessageModule = &Module{}
)

type Module struct {
	log            *zerolog.Logger
	broker         *broker.Broker
	client         *grpcClient.Client
	cdc            codec.Codec
	parser         messages.MessageAddressesParser
	enabledModules []string // xxx fixme
}

func New(b *broker.Broker, cli *grpcClient.Client, cdc codec.Codec, parser messages.MessageAddressesParser,
	modules []string) *Module {

	l := zerolog.New(os.Stderr).Output(zerolog.ConsoleWriter{Out: os.Stderr}).With().Timestamp().
		Str("module", "staking").Logger()

	return &Module{
		log:            &l,
		broker:         b,
		client:         cli,
		cdc:            cdc,
		parser:         parser,
		enabledModules: modules,
	}
}

func (m *Module) Name() string { return "staking" }
