package broker

import (
	"context"
)

func (b *Broker) PublishRawBlock(_ context.Context, block interface{}) error {
	return b.marshalAndProduce(RawBlock, block)
}

func (b *Broker) PublishRawTransaction(_ context.Context, tx interface{}) error {
	return b.marshalAndProduce(RawTransaction, tx)
}

func (b *Broker) PublishRawBlockResults(_ context.Context, br interface{}) error {
	return b.marshalAndProduce(RawBlockResults, br)
}
