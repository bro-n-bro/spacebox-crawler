package broker

import (
	"context"

	"github.com/bro-n-bro/spacebox/broker/model"
)

func (b *Broker) PublishGrantMessage(ctx context.Context, gm model.GrantMessage) error {
	return b.marshalAndProduce(GrantMessage, gm)
}
