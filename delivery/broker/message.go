package broker

import (
	"context"

	"github.com/bro-n-bro/spacebox/broker/model"
)

func (b *Broker) PublishMessage(ctx context.Context, message model.Message) error {
	return b.marshalAndProduce(Message, message)
}
