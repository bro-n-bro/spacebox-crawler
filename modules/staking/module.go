package staking

import (
	"bro-n-bro-osmosis/internal/rep"
	tb "bro-n-bro-osmosis/pkg/mapper/to_broker"
	"bro-n-bro-osmosis/types"
	"os"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/rs/zerolog"

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
	client         *grpcClient.Client
	broker         rep.Broker
	tbM            tb.ToBroker
	cdc            codec.Codec
	parser         messages.MessageAddressesParser
	enabledModules []string // xxx fixme
}

func New(b rep.Broker, cli *grpcClient.Client, tbM tb.ToBroker, cdc codec.Codec,
	modules []string) *Module {

	l := zerolog.New(os.Stderr).Output(zerolog.ConsoleWriter{Out: os.Stderr}).With().Timestamp().
		Str("module", "staking").Logger()

	return &Module{
		log:            &l,
		broker:         b,
		client:         cli,
		tbM:            tbM,
		cdc:            cdc,
		enabledModules: modules,
	}
}

func (m *Module) Name() string { return "staking" }
