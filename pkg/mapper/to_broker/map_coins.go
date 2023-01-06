package tobroker

import (
	"github.com/hexy-dev/spacebox/broker/model"

	"github.com/hexy-dev/spacebox-crawler/types"
)

func (tb ToBroker) MapCoin(coin types.Coin) model.Coin {
	return model.Coin{
		Denom:  coin.Denom,
		Amount: coin.Amount,
	}
}

func (tb ToBroker) MapCoins(coins types.Coins) model.Coins {
	res := make(model.Coins, len(coins))
	for i, c := range coins {
		res[i] = tb.MapCoin(c)
	}
	return res
}
