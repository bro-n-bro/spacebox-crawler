package broker

import (
	"context"

	"github.com/bro-n-bro/spacebox/broker/model"
)

func (b *Broker) PublishSupply(ctx context.Context, supply model.Supply) error {
	return b.marshalAndProduce(Supply, supply)
}
