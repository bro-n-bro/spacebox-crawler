package broker

import (
	"context"

	"github.com/bro-n-bro/spacebox/broker/model"
)

func (b *Broker) PublishRedelegationMessage(ctx context.Context, rm model.RedelegationMessage) error {
	return b.marshalAndProduce(RedelegationMessage, rm)
}
