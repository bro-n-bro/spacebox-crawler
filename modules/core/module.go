package core

import (
	"os"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/rs/zerolog"

	"github.com/bro-n-bro/spacebox-crawler/internal/rep"
	tb "github.com/bro-n-bro/spacebox-crawler/pkg/mapper/to_broker"
	"github.com/bro-n-bro/spacebox-crawler/types"
)

var (
	_ types.Module             = &Module{}
	_ types.BlockHandler       = &Module{}
	_ types.MessageHandler     = &Module{}
	_ types.TransactionHandler = &Module{}
	_ types.ValidatorsHandler  = &Module{}
)

type Module struct {
	log    *zerolog.Logger
	tbM    tb.ToBroker
	broker broker
	cdc    codec.Codec
	parser MessageAddressesParser
}

func New(b rep.Broker, tbM tb.ToBroker, cdc codec.Codec, parser MessageAddressesParser) *Module {
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
