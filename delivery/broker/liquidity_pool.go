package broker

import (
	"context"

	"github.com/bro-n-bro/spacebox/broker/model"
)

func (b *Broker) PublishLiquidityPool(_ context.Context, v model.LiquidityPool) error {
	return b.marshalAndProduce(LiquidityPool, v)
}
