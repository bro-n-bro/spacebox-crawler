package liquidity

import (
	"context"

	"github.com/bro-n-bro/spacebox/broker/model"
)

type broker interface {
	PublishSwap(context.Context, model.Swap) error
	PublishLiquidityPool(context.Context, model.LiquidityPool) error
}
