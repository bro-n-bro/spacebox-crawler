package broker

import (
	"context"

	"github.com/bro-n-bro/spacebox/broker/model"
)

func (b *Broker) PublishMultiSendMessage(ctx context.Context, msm model.MultiSendMessage) error {
	return b.marshalAndProduce(MultiSendMessage, msm)
}
