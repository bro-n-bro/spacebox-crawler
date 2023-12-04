package broker

import (
	"context"

	"github.com/bro-n-bro/spacebox/broker/model"
)

func (b *Broker) PublishDistributionCommission(_ context.Context, commission model.DistributionCommission) error {
	return b.marshalAndProduce(DistributionCommission, commission)
}
