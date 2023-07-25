package liquidity

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"gopkg.in/yaml.v3"
)

const (
	// MinReserveCoinNum is the minimum number of reserve coins in each liquidity pool.
	MinReserveCoinNum uint32 = 2

	// MaxReserveCoinNum is the maximum number of reserve coins in each liquidity pool.
	MaxReserveCoinNum uint32 = 2
)

var (
	MinOfferCoinAmount = sdk.NewInt(100)
)

// String returns a human readable string representation of the parameters.
func (p Params) String() string {
	out, _ := yaml.Marshal(p)
	return string(out)
}
