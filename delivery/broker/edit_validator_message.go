package broker

import (
	"context"

	"github.com/bro-n-bro/spacebox/broker/model"
)

func (b *Broker) PublishEditValidatorMessage(_ context.Context, msg model.EditValidatorMessage) error {
	return b.marshalAndProduce(EditValidatorMessage, msg)
}
