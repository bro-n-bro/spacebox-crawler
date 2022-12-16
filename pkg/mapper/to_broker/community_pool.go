package to_broker

import (
	"github.com/hexy-dev/spacebox/broker/model"

	"bro-n-bro-osmosis/types"
)

func (tb ToBroker) MapCommunityPool(cp types.CommunityPool) model.CommunityPool {
	return model.CommunityPool{
		Coins:  tb.MapCoins(cp.Coins),
		Height: cp.Height,
	}
}
