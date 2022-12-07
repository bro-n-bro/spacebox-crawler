package to_broker

import "github.com/cosmos/cosmos-sdk/codec"

// ToBroker mapper
type ToBroker struct {
	cdc codec.Codec
}

func NewToBroker(cdc codec.Codec) *ToBroker {
	return &ToBroker{cdc: cdc}
}
