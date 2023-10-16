package broker

import (
	"context"

	"github.com/bro-n-bro/spacebox/broker/model"
)

func (b *Broker) PublishDelegation(ctx context.Context, d model.Delegation) error {
	d.IsActive = true
	return b.marshalAndProduce(Delegation, d)
}

func (b *Broker) PublishDisabledDelegation(ctx context.Context, d model.Delegation) error {
	d.IsActive = false
	return b.marshalAndProduce(Delegation, d)
}
