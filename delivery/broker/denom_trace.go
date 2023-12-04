package broker

import (
	"context"

	"github.com/bro-n-bro/spacebox/broker/model"
)

func (b *Broker) PublishDenomTrace(ctx context.Context, dt model.DenomTrace) error {
	return b.marshalAndProduce(DenomTrace, dt)
}
