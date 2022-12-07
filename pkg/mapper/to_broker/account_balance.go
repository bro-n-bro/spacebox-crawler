package to_broker

import (
	"bro-n-bro-osmosis/adapter/broker/model"
	"bro-n-bro-osmosis/types"
)

func (tb ToBroker) MapAccountBalance(ab types.AccountBalance) model.AccountBalance {
	return model.AccountBalance{
		Address: ab.Address,
		Height:  ab.Height,
		Coins:   tb.MapCoins(ab.Balance),
	}
}
