package broker

import (
	"context"

	"github.com/bro-n-bro/spacebox/broker/model"
)

func (b *Broker) PublishUnbondingDelegation(ctx context.Context, ud model.UnbondingDelegation) error {
	return b.marshalAndProduce(UnbondingDelegation, ud)
}
