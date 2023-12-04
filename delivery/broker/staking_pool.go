package broker

import (
	"context"

	"github.com/bro-n-bro/spacebox/broker/model"
)

func (b *Broker) PublishStakingPool(ctx context.Context, sp model.StakingPool) error {
	return b.marshalAndProduce(StakingPool, sp)
}
