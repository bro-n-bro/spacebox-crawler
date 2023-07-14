package auth

import (
	"os"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/rs/zerolog"

	grpcClient "github.com/bro-n-bro/spacebox-crawler/client/grpc"
	"github.com/bro-n-bro/spacebox-crawler/modules/core"
	tb "github.com/bro-n-bro/spacebox-crawler/pkg/mapper/to_broker"
	"github.com/bro-n-bro/spacebox-crawler/types"
)

const (
	moduleName = "auth"
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

func New(b broker, cli *grpcClient.Client, tb tb.ToBroker, cdc codec.Codec,
	parser core.MessageAddressesParser) *Module {

	l := zerolog.New(os.Stderr).Output(zerolog.ConsoleWriter{Out: os.Stderr}).With().Timestamp().
		Str("module", moduleName).Logger()

	return &Module{
		log:    &l,
		broker: b,
		tbM:    tb,
		client: cli,
		cdc:    cdc,
		parser: parser,
	}
}

func (m *Module) Name() string { return moduleName }
