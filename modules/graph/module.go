package graph

import (
	"os"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/rs/zerolog"

	tb "github.com/bro-n-bro/spacebox-crawler/pkg/mapper/to_broker"
	"github.com/bro-n-bro/spacebox-crawler/types"
)

const (
	keyModule  = "module"
	moduleName = "graph"
)

var (
	_ types.Module         = &Module{}
	_ types.MessageHandler = &Module{}
)

type (
	Module struct {
		log    *zerolog.Logger
		broker broker
		tbM    tb.ToBroker
		cdc    codec.Codec
	}
)

func New(b broker, tbM tb.ToBroker, cdc codec.Codec) *Module {
	l := zerolog.New(os.Stderr).
		Output(zerolog.ConsoleWriter{Out: os.Stderr}).
		With().Timestamp().Str(keyModule, moduleName).Logger()

	return &Module{
		log:    &l,
		broker: b,
		tbM:    tbM,
		cdc:    cdc,
	}
}

func (m *Module) Name() string { return moduleName }
