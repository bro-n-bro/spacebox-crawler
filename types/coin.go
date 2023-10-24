package types

import (
	"math"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type (
	Coins []Coin

	Coin struct {
		Denom  string  `json:"denom"`
		Amount float64 `json:"amount"`
	}
)

func NewCoin(denom string, amount float64) Coin {
	return Coin{
		Denom:  denom,
		Amount: amount,
	}
}

func NewCoinFromSDK(coin sdk.Coin) Coin {
	return Coin{
		Denom:  coin.Denom,
		Amount: float64(coin.Amount.BigInt().Int64()),
	}
}

func NewCoinsFromSDK(coins sdk.Coins) Coins {
	res := make(Coins, len(coins))
	for i, c := range coins {
		res[i] = NewCoinFromSDK(c)
	}

	return res
}

func NewCoinsFromSDKDec(coins sdk.DecCoins) Coins {
	res := make(Coins, len(coins))
	for i, c := range coins {
		res[i] = Coin{
			Denom:  c.Denom,
			Amount: c.Amount.MustFloat64(),
		}
	}

	return res
}

func (c Coin) IsEqual(o Coin) bool {
	return c.Denom == o.Denom && math.Nextafter(c.Amount, o.Amount) == o.Amount
}
