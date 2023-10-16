package broker

import (
	"context"

	"github.com/bro-n-bro/spacebox/broker/model"
)

func (b *Broker) PublishRankParams(ctx context.Context, msg model.RankParams) error {
	return b.marshalAndProduce(RankParams, msg)
}
