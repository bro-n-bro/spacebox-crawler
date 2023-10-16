package broker

import (
	"context"

	"github.com/bro-n-bro/spacebox/broker/model"
)

func (b *Broker) PublishGridParams(ctx context.Context, msg model.GridParams) error {
	return b.marshalAndProduce(GridParams, msg)
}
