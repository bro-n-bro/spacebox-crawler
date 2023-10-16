package broker

import (
	"context"

	"github.com/bro-n-bro/spacebox/broker/model"
)

func (b *Broker) PublishReceivePacketMessage(ctx context.Context, r model.RecvPacketMessage) error {
	return b.marshalAndProduce(ReceivePacketMessage, r)
}
