package broker

import (
	"context"

	"github.com/bro-n-bro/spacebox/broker/model"
)

func (b *Broker) PublishAccount(ctx context.Context, account model.Account) error {
	return b.marshalAndProduce(Account, account)
}
