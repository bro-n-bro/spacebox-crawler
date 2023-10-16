package broker

import (
	"context"

	"github.com/bro-n-bro/spacebox/broker/model"
)

func (b *Broker) PublishValidatorPreCommit(_ context.Context, v model.ValidatorPreCommit) error {
	return b.marshalAndProduce(ValidatorPreCommit, v)
}
