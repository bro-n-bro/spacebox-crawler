package broker

import (
	"context"

	"github.com/bro-n-bro/spacebox/broker/model"
)

func (b *Broker) PublishProposerReward(ctx context.Context, r model.ProposerReward) error {
	return b.marshalAndProduce(ProposerReward, r)
}
