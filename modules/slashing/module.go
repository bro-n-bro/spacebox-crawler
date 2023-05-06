package slashing

import (
	"os"

	"github.com/rs/zerolog"

	"github.com/bro-n-bro/spacebox-crawler/types"
)

const (
	moduleName = "slashing"
)

var (
	_ types.Module         = &Module{}
	_ types.MessageHandler = &Module{}
)

type Module struct {
	log    *zerolog.Logger
	broker broker
}

func New(b broker) *Module {
	l := zerolog.New(os.Stderr).Output(zerolog.ConsoleWriter{Out: os.Stderr}).With().Timestamp().
		Str("module", moduleName).Logger()

	return &Module{
		log:    &l,
		broker: b,
	}
}

func (m *Module) Name() string { return moduleName }
