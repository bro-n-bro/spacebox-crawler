package broker

import (
	"context"

	"github.com/bro-n-bro/spacebox/broker/model"
)

func (b *Broker) PublishUnjailMessage(_ context.Context, msg model.UnjailMessage) error {
	return b.marshalAndProduce(UnjailMessage, msg)
}
