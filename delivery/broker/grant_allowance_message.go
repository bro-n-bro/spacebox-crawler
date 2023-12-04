package broker

import (
	"context"

	"github.com/bro-n-bro/spacebox/broker/model"
)

func (b *Broker) PublishGrantAllowanceMessage(ctx context.Context, gam model.GrantAllowanceMessage) error {
	return b.marshalAndProduce(GrantAllowanceMessage, gam)
}
