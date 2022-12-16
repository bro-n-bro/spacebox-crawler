package to_broker

import (
	"github.com/hexy-dev/spacebox/broker/model"

	"bro-n-bro-osmosis/types"
)

func (tb ToBroker) MapRedelegation(m types.Redelegation) model.Redelegation {
	return model.Redelegation{
		CompletionTime:   m.CompletionTime,
		Coin:             tb.MapCoin(m.Coin),
		DelegatorAddress: m.DelegatorAddress,
		SrcValidator:     m.SrcValidator,
		DstValidator:     m.DstValidator,
		Height:           m.Height,
	}
}
