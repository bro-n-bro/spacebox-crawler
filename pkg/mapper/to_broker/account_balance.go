package to_broker

import (
	"github.com/hexy-dev/spacebox-crawler/types"

	"github.com/hexy-dev/spacebox/broker/model"
)

func (tb ToBroker) MapAccountBalance(ab types.AccountBalance) model.AccountBalance {
	return model.AccountBalance{
		Address: ab.Address,
		Height:  ab.Height,
		Coins:   tb.MapCoins(ab.Balance),
	}
}
