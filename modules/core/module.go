package core

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/rs/zerolog"

	"github.com/bro-n-bro/spacebox-crawler/modules/utils"
	tb "github.com/bro-n-bro/spacebox-crawler/pkg/mapper/to_broker"
	"github.com/bro-n-bro/spacebox-crawler/types"
)

const (
	ModuleName = "core"
)

var (
	_ types.Module             = &Module{}
	_ types.BlockHandler       = &Module{}
	_ types.MessageHandler     = &Module{}
	_ types.TransactionHandler = &Module{}
)

type Module struct {
	log    *zerolog.Logger
	tbM    tb.ToBroker
	broker broker
	cdc    codec.Codec
	parser MsgAddrParser
}

func New(b broker, tbM tb.ToBroker, cdc codec.Codec, parser MsgAddrParser) *Module {
	return &Module{
		log:    utils.NewModuleLogger(ModuleName),
		broker: b,
		tbM:    tbM,
		cdc:    cdc,
		parser: parser,
	}
}

func (m *Module) Name() string { return ModuleName }
