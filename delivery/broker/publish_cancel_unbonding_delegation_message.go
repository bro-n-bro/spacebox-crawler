package broker

import (
	"context"

	"github.com/bro-n-bro/spacebox/broker/model"
)

func (b *Broker) PublishCancelUnbondingDelegationMessage(_ context.Context, d model.CancelUnbondingDelegationMessage) error {
	return b.marshalAndProduce(CancelUnbondingDelegationMessage, d)
}
