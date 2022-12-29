package to_broker

import (
	"github.com/hexy-dev/spacebox/broker/model"

	"github.com/hexy-dev/spacebox-crawler/types"
)

func (tb ToBroker) MapDelegationRewardMessage(m types.DelegationRewardMessage) model.DelegationRewardMessage {
	return model.DelegationRewardMessage{
		Coins:            tb.MapCoins(m.Coins),
		DelegatorAddress: m.DelegatorAddress,
		ValidatorAddress: m.ValidatorAddress,
		TxHash:           m.TxHash,
		Height:           m.Height,
	}
}
