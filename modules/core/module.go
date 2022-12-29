package core

import (
	"os"

	"github.com/hexy-dev/spacebox-crawler/internal/rep"
	"github.com/hexy-dev/spacebox-crawler/modules/messages"
	tb "github.com/hexy-dev/spacebox-crawler/pkg/mapper/to_broker"
	"github.com/hexy-dev/spacebox-crawler/types"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/rs/zerolog"
)

var (
	_ types.Module            = &Module{}
	_ types.BlockModule       = &Module{}
	_ types.MessageModule     = &Module{}
	_ types.TransactionModule = &Module{}
)

type Module struct {
	log    *zerolog.Logger
	tbM    tb.ToBroker
	broker rep.Broker
	cdc    codec.Codec
	parser messages.MessageAddressesParser
}

func New(b rep.Broker, tbM tb.ToBroker, cdc codec.Codec, parser messages.MessageAddressesParser) *Module {
	l := zerolog.New(os.Stderr).Output(zerolog.ConsoleWriter{Out: os.Stderr}).With().Timestamp().
		Str("module", "core").Logger()

	return &Module{
		log:    &l,
		broker: b,
		tbM:    tbM,
		cdc:    cdc,
		parser: parser,
	}
}

func (m *Module) Name() string { return "core" }
