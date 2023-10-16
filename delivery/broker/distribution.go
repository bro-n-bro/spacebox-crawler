package broker

import (
	"context"

	"github.com/bro-n-bro/spacebox/broker/model"
)

func (b *Broker) PublishDistributionParams(ctx context.Context, dp model.DistributionParams) error {
	return b.marshalAndProduce(DistributionParams, dp)
}
