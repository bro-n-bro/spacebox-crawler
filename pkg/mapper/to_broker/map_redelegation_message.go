package to_broker

import (
	"github.com/hexy-dev/spacebox-crawler/types"
	"github.com/hexy-dev/spacebox/broker/model"
)

func (tb ToBroker) MapRedelegationMessage(m types.RedelegationMessage) model.RedelegationMessage {
	return model.RedelegationMessage{
		Redelegation: tb.MapRedelegation(m.Redelegation),
		TxHash:       m.TxHash,
	}
}
