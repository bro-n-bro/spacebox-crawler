package to_broker

import (
	"github.com/hexy-dev/spacebox/broker/model"

	"bro-n-bro-osmosis/types"
)

func (tb ToBroker) MapSupply(supply types.TotalSupply) model.Supply {
	return model.Supply{
		Height: supply.Height,
		Coins:  tb.MapCoins(supply.Coins),
	}
}
