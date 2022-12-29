package distribution

import (
	"os"

	grpcClient "github.com/hexy-dev/spacebox-crawler/client/grpc"
	"github.com/hexy-dev/spacebox-crawler/internal/rep"
	"github.com/hexy-dev/spacebox-crawler/modules/messages"
	tb "github.com/hexy-dev/spacebox-crawler/pkg/mapper/to_broker"
	"github.com/hexy-dev/spacebox-crawler/types"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/rs/zerolog"
)

var (
	_ types.Module        = &Module{}
	_ types.BlockModule   = &Module{}
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

func New(b rep.Broker, cli *grpcClient.Client, tbM tb.ToBroker, cdc codec.Codec,
	parser messages.MessageAddressesParser) *Module {

	l := zerolog.New(os.Stderr).Output(zerolog.ConsoleWriter{Out: os.Stderr}).With().Timestamp().
		Str("module", "distribution").Logger()

	return &Module{
		log:    &l,
		broker: b,
		client: cli,
		tbM:    tbM,
		cdc:    cdc,
		parser: parser,
	}
}

func (m *Module) Name() string { return "distribution" }
