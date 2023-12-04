package broker

import (
	"context"

	"github.com/bro-n-bro/spacebox/broker/model"
)

func (b *Broker) PublishDelegationMessage(ctx context.Context, dm model.DelegationMessage) error {
	return b.marshalAndProduce(DelegationMessage, dm)
}
