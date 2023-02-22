package liquidity

import sdk "github.com/cosmos/cosmos-sdk/types"

const (
	// MinReserveCoinNum is the minimum number of reserve coins in each liquidity pool.
	MinReserveCoinNum uint32 = 2

	// MaxReserveCoinNum is the maximum number of reserve coins in each liquidity pool.
	MaxReserveCoinNum uint32 = 2
)

var (
	MinOfferCoinAmount = sdk.NewInt(100)
)
