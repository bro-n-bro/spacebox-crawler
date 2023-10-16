package broker

import (
	"context"

	"github.com/bro-n-bro/spacebox/broker/model"
)

func (b *Broker) PublishValidatorDescription(ctx context.Context, description model.ValidatorDescription) error {
	if checkCache(description.OperatorAddress, description.Height, b.cache.valDescription) {
		return b.marshalAndProduce(ValidatorDescription, description)
	}

	return nil
}
