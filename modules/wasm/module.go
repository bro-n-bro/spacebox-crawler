package wasm

import (
	"github.com/rs/zerolog"

	"github.com/bro-n-bro/spacebox-crawler/modules/utils"
	"github.com/bro-n-bro/spacebox-crawler/types"
)

const (
	ModuleName = "wasm"
)

var (
	_ types.Module         = &Module{}
	_ types.MessageHandler = &Module{}
)

type (
	Module struct {
		log    *zerolog.Logger
		broker broker
	}
)

func New(b broker) *Module {
	return &Module{
		log:    utils.NewModuleLogger(ModuleName),
		broker: b,
	}
}

func (m *Module) Name() string { return ModuleName }
