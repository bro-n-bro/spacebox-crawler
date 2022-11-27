package to_broker

import (
	"bro-n-bro-osmosis/adapter/broker/model"
	"bro-n-bro-osmosis/types"
)

func (tb ToBroker) MapDelegationRewardMessage(m types.DelegationRewardMessage) model.DelegationRewardMessage {
	return model.DelegationRewardMessage{
		Coins:            tb.MapCoins(m.Coins),
		DelegatorAddress: m.DelegatorAddress,
		ValidatorAddress: m.ValidatorAddress,
	}
}
