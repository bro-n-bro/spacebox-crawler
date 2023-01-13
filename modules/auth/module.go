package auth

import (
	"os"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/rs/zerolog"

	grpcClient "github.com/hexy-dev/spacebox-crawler/client/grpc"
	"github.com/hexy-dev/spacebox-crawler/internal/rep"
	"github.com/hexy-dev/spacebox-crawler/modules/core"
	tb "github.com/hexy-dev/spacebox-crawler/pkg/mapper/to_broker"
	"github.com/hexy-dev/spacebox-crawler/types"
)

var (
	_ types.Module         = &Module{}
	_ types.GenesisHandler = &Module{}
	_ types.MessageHandler = &Module{}
)

type Module struct {
	log    *zerolog.Logger
	client *grpcClient.Client
	broker broker
	tbM    tb.ToBroker
	cdc    codec.Codec
	parser core.MessageAddressesParser
}

func New(b rep.Broker, cli *grpcClient.Client, tb tb.ToBroker, cdc codec.Codec,
	parser core.MessageAddressesParser) *Module {

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
