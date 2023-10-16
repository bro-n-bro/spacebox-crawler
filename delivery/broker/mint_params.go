package broker

import (
	"context"

	"github.com/bro-n-bro/spacebox/broker/model"
)

func (b *Broker) PublishMintParams(ctx context.Context, mp model.MintParams) error {
	return b.marshalAndProduce(MintParams, mp)
}
