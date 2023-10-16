package broker

import (
	"context"

	"github.com/bro-n-bro/spacebox/broker/model"
)

func (b *Broker) PublishRedelegation(ctx context.Context, r model.Redelegation) error {
	return b.marshalAndProduce(Redelegation, r)
}
