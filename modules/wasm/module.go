package wasm

import (
	"github.com/cosmos/cosmos-sdk/codec"
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
		cdc    codec.Codec
	}
)

func New(b broker, cdc codec.Codec) *Module {
	return &Module{
		log:    utils.NewModuleLogger(ModuleName),
		broker: b,
		cdc:    cdc,
	}
}

func (m *Module) Name() string { return ModuleName }
