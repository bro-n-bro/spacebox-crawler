package broker

import (
	"context"

	"github.com/bro-n-bro/spacebox/broker/model"
)

func (b *Broker) PublishValidatorCommission(ctx context.Context, commission model.ValidatorCommission) error {
	if checkCache(commission.OperatorAddress, commission.Height, b.cache.valCommission) {
		return b.marshalAndProduce(ValidatorCommission, commission)
	}

	return nil
}
