package broker

import (
	"context"

	"github.com/bro-n-bro/spacebox/broker/model"
)

func (b *Broker) PublishFeeAllowance(ctx context.Context, feeAllowance model.FeeAllowance) error {
	return b.marshalAndProduce(FeeAllowance, feeAllowance)
}
