package staking

import (
	"os"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/rs/zerolog"

	grpcClient "github.com/bro-n-bro/spacebox-crawler/client/grpc"
	"github.com/bro-n-bro/spacebox-crawler/internal/rep"
	tb "github.com/bro-n-bro/spacebox-crawler/pkg/mapper/to_broker"
	"github.com/bro-n-bro/spacebox-crawler/types"
)

var (
	_ types.Module         = &Module{}
	_ types.GenesisHandler = &Module{}
	_ types.BlockHandler   = &Module{}
	_ types.MessageHandler = &Module{}
)

type Module struct {
	log            *zerolog.Logger
	client         *grpcClient.Client
	broker         broker
	tbM            tb.ToBroker
	cdc            codec.Codec
	enabledModules []string // xxx fixme
}

func New(b rep.Broker, cli *grpcClient.Client, tbM tb.ToBroker, cdc codec.Codec, modules []string) *Module {
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
