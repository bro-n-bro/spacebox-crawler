package broker

import (
	"context"

	"github.com/bro-n-bro/spacebox/broker/model"
)

func (b *Broker) PublishGovParams(ctx context.Context, params model.GovParams) error {
	return b.marshalAndProduce(GovParams, params)
}
