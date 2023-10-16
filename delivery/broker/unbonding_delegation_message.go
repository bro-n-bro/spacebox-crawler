package broker

import (
	"context"

	"github.com/bro-n-bro/spacebox/broker/model"
)

func (b *Broker) PublishUnbondingDelegationMessage(ctx context.Context, udm model.UnbondingDelegationMessage) error {
	return b.marshalAndProduce(UnbondingDelegationMessage, udm)
}
