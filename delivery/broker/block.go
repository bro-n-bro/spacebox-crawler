package broker

import (
	"context"

	"github.com/bro-n-bro/spacebox/broker/model"
)

func (b *Broker) PublishBlock(ctx context.Context, block model.Block) error {
	return b.marshalAndProduce(Block, block)
}
