package to_broker

import (
	"github.com/hexy-dev/spacebox/broker/model"

	"bro-n-bro-osmosis/types"
)

func (tb ToBroker) MapRedelegationMessage(m types.RedelegationMessage) model.RedelegationMessage {
	return model.RedelegationMessage{
		Redelegation: tb.MapRedelegation(m.Redelegation),
		TxHash:       m.TxHash,
	}
}
