package to_broker

import (
	"github.com/hexy-dev/spacebox/broker/model"

	"github.com/hexy-dev/spacebox-crawler/types"
)

func (tb ToBroker) MapCommunityPool(cp types.CommunityPool) model.CommunityPool {
	return model.CommunityPool{
		Coins:  tb.MapCoins(cp.Coins),
		Height: cp.Height,
	}
}
