package broker

import (
	"context"

	"github.com/bro-n-bro/spacebox/broker/model"
)

func (b *Broker) PublishSwap(ctx context.Context, swap model.Swap) error {
	return b.marshalAndProduce(Swap, swap)
}
