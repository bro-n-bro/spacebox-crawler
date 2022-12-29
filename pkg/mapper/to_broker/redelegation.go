package to_broker

import (
	"github.com/hexy-dev/spacebox/broker/model"

	"github.com/hexy-dev/spacebox-crawler/types"
)

func (tb ToBroker) MapRedelegation(m types.Redelegation) model.Redelegation {
	return model.Redelegation{
		CompletionTime:      m.CompletionTime,
		Coin:                tb.MapCoin(m.Coin),
		DelegatorAddress:    m.DelegatorAddress,
		SrcValidatorAddress: m.SrcValidator,
		DstValidatorAddress: m.DstValidator,
		Height:              m.Height,
	}
}
