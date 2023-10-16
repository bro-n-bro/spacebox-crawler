package tobroker

import "github.com/cosmos/cosmos-sdk/codec"

type (
	// ToBroker mapper
	ToBroker struct {
		cdc   codec.Codec
		amino *codec.LegacyAmino
	}
)

func NewToBroker(cdc codec.Codec, amino *codec.LegacyAmino) *ToBroker {
	return &ToBroker{
		cdc:   cdc,
		amino: amino,
	}
}
