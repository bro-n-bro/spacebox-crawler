package broker

import (
	"context"

	"github.com/bro-n-bro/spacebox/broker/model"
)

func (b *Broker) PublishCommunityPool(ctx context.Context, cp model.CommunityPool) error {
	return b.marshalAndProduce(CommunityPool, cp)
}
