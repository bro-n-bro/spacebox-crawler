package broker

import (
	"context"

	"github.com/bro-n-bro/spacebox/broker/model"
)

func (b *Broker) PublishUnJailMessage(_ context.Context, msg model.UnjailMessage) error {
	return b.marshalAndProduce(UnJailMessage, msg)
}
