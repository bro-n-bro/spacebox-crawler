package types

import (
	"math"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Coins []Coin

type Coin struct {
	Denom  string  `json:"denom"`
	Amount float64 `json:"amount"`
}

func NewCoinFromCdk(coin sdk.Coin) Coin {
	return Coin{
		Denom:  coin.Denom,
		Amount: float64(coin.Amount.BigInt().Int64()),
	}
}

func NewCoinsFromCdk(coins sdk.Coins) Coins {
	res := make(Coins, len(coins))
	for i, c := range coins {
		res[i] = NewCoinFromCdk(c)
	}
	return res
}

func NewCoinsFromCdkDec(coins sdk.DecCoins) Coins {
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
