package to_broker

import (
	"bro-n-bro-osmosis/adapter/broker/model"
	"bro-n-bro-osmosis/types"
)

func (tb ToBroker) MapUnbondingDelegationMessage(udm types.UnbondingDelegationMessage) model.UnbondingDelegationMessage {
	return model.UnbondingDelegationMessage{
		CompletionTimestamp: udm.CompletionTimestamp,
		Coin:                tb.MapCoin(udm.Coin),
		DelegatorAddress:    udm.DelegatorAddress,
		ValidatorOperAddr:   udm.ValidatorOperAddr,
		TxHash:              udm.TxHash,
		Height:              udm.Height,
	}
}
