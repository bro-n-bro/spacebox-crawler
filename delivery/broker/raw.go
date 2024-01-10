package broker

import (
	"context"
)

func (b *Broker) PublishRawBlock(_ context.Context, block interface{}) error {
	return b.marshalAndProduce(Account, block)
}

func (b *Broker) PublishRawTransaction(_ context.Context, tx interface{}) error {
	return b.marshalAndProduce(Account, tx)
}

func (b *Broker) PublishRawBlockResults(_ context.Context, br interface{}) error {
	return b.marshalAndProduce(Account, br)
}
