package broker

import (
	"context"

	"github.com/bro-n-bro/spacebox/broker/model"
)

func (b *Broker) PublishValidator(ctx context.Context, val model.Validator) error {
	if checkCache(val.ConsensusAddress, val.Height, b.cache.validator) {
		return b.marshalAndProduce(Validator, val)
	}

	return nil
}
