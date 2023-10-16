package broker

import (
	"context"

	"github.com/bro-n-bro/spacebox/broker/model"
)

func (b *Broker) PublishCreateValidatorMessage(_ context.Context, cvm model.CreateValidatorMessage) error {
	return b.marshalAndProduce(CreateValidatorMessage, cvm)
}
