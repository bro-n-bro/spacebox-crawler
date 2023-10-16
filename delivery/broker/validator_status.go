package broker

import (
	"context"

	"github.com/bro-n-bro/spacebox/broker/model"
)

func (b *Broker) PublishValidatorStatus(ctx context.Context, status model.ValidatorStatus) error {
	if checkCache(status.ConsensusAddress, status.Height, b.cache.valStatus) {
		return b.marshalAndProduce(ValidatorStatus, status)
	}

	return nil
}
