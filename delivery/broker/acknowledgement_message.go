package broker

import (
	"context"

	"github.com/bro-n-bro/spacebox/broker/model"
)

func (b *Broker) PublishAcknowledgementMessage(ctx context.Context, msg model.AcknowledgementMessage) error {
	return b.marshalAndProduce(AcknowledgementMessage, msg)
}
