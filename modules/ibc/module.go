package ibc

import (
	"os"

	"github.com/rs/zerolog"

	tb "github.com/bro-n-bro/spacebox-crawler/pkg/mapper/to_broker"
	"github.com/bro-n-bro/spacebox-crawler/types"
)

const (
	moduleName = "ibc"
)

var (
	_ types.Module         = &Module{}
	_ types.MessageHandler = &Module{}
)

type Module struct {
	log    *zerolog.Logger
	broker broker
	tbM    tb.ToBroker
}

func New(b broker, tbM tb.ToBroker) *Module {
	l := zerolog.New(os.Stderr).Output(zerolog.ConsoleWriter{Out: os.Stderr}).With().Timestamp().
		Str("module", moduleName).Logger()

	return &Module{
		log:    &l,
		broker: b,
		tbM:    tbM,
	}
}

func (m *Module) Name() string { return moduleName }
