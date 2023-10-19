package broker

import (
	"context"

	"github.com/bro-n-bro/spacebox/broker/model"
)

func (b *Broker) PublishInvestmintMessage(ctx context.Context, msg model.InvestmintMessage) error {
	return b.marshalAndProduce(InvestmintMessage, msg)
}
