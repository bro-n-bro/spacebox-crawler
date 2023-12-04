package broker

import (
	"context"

	"github.com/bro-n-bro/spacebox/broker/model"
)

func (b *Broker) PublishAccountBalance(ctx context.Context, ab model.AccountBalance) error {
	return b.marshalAndProduce(AccountBalance, ab)
}
