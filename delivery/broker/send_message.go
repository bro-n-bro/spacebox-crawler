package broker

import (
	"context"

	"github.com/bro-n-bro/spacebox/broker/model"
)

func (b *Broker) PublishSendMessage(ctx context.Context, sm model.SendMessage) error {
	return b.marshalAndProduce(SendMessage, sm)
}
