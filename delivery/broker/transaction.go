package broker

import (
	"context"

	"github.com/bro-n-bro/spacebox/broker/model"
)

func (b *Broker) PublishTransaction(ctx context.Context, tx model.Transaction) error {
	return b.marshalAndProduce(Transaction, tx)
}
