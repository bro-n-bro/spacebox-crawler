package tobroker

import "github.com/cosmos/cosmos-sdk/codec"

// ToBroker mapper
type ToBroker struct {
	cdc   codec.Codec
	amino *codec.LegacyAmino
}

func NewToBroker(cdc codec.Codec, amino *codec.LegacyAmino) *ToBroker {
	return &ToBroker{cdc: cdc, amino: amino}
}
