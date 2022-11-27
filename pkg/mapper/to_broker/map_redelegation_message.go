package to_broker

import (
	"bro-n-bro-osmosis/adapter/broker/model"
	"bro-n-bro-osmosis/types"
)

func (tb ToBroker) MapRedelegationMessage(m types.RedelegationMessage) model.RedelegationMessage {
	return model.RedelegationMessage{
		CompletionTime:   m.CompletionTime,
		Coin:             tb.MapCoin(m.Coin),
		DelegatorAddress: m.DelegatorAddress,
		SrcValidator:     m.SrcValidator,
		DstValidator:     m.DstValidator,
		TxHash:           m.TxHash,
		Height:           m.Height,
	}
}
