package broker

import (
	"context"

	"github.com/bro-n-bro/spacebox/broker/model"
)

func (b *Broker) PublishDMNParams(ctx context.Context, msg model.DMNParams) error {
	return b.marshalAndProduce(DMNParams, msg)
}
