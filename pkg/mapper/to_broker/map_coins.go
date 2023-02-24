package tobroker

import (
	jsoniter "github.com/json-iterator/go"

	"github.com/bro-n-bro/spacebox-crawler/types"
	"github.com/bro-n-bro/spacebox/broker/model"
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

func (tb ToBroker) MapCoinToString(coin types.Coin) string {
	coinsStr, _ := jsoniter.MarshalToString(&coin)
	return coinsStr
}

func (tb ToBroker) MapCoinsToString(coins types.Coins) string {
	res := make(model.Coins, len(coins))
	for i, c := range coins {
		res[i] = tb.MapCoin(c)
	}

	coinsStr, _ := jsoniter.MarshalToString(&res)

	return coinsStr
}
