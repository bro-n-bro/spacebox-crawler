package broker

import (
	"context"

	"github.com/bro-n-bro/spacebox/broker/model"
)

func (b *Broker) PublishDelegationReward(ctx context.Context, dr model.DelegationReward) error {
	return b.marshalAndProduce(DelegationReward, dr)
}
