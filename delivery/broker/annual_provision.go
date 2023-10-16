package broker

import (
	"context"

	"github.com/bro-n-bro/spacebox/broker/model"
)

func (b *Broker) PublishAnnualProvision(ctx context.Context, ap model.AnnualProvision) error {
	return b.marshalAndProduce(AnnualProvision, ap)
}
