package broker

import (
	"context"

	"github.com/bro-n-bro/spacebox/broker/model"
)

func (b *Broker) PublishStakingParams(ctx context.Context, sp model.StakingParams) error {
	return b.marshalAndProduce(StakingParams, sp)
}
