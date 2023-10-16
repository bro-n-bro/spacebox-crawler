package broker

import (
	"context"

	"github.com/bro-n-bro/spacebox/broker/model"
)

func (b *Broker) PublishTransferMessage(ctx context.Context, tm model.TransferMessage) error {
	return b.marshalAndProduce(TransferMessage, tm)
}
