package broker

import (
	"context"

	"github.com/bro-n-bro/spacebox/broker/model"
)

func (b *Broker) PublishRevokeAllowanceMessage(ctx context.Context, ram model.RevokeAllowanceMessage) error {
	return b.marshalAndProduce(RevokeAllowanceMessage, ram)
}
