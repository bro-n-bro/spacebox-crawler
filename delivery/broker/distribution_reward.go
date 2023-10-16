package broker

import (
	"context"

	"github.com/bro-n-bro/spacebox/broker/model"
)

func (b *Broker) PublishDistributionReward(_ context.Context, reward model.DistributionReward) error {
	return b.marshalAndProduce(DistributionReward, reward)
}
