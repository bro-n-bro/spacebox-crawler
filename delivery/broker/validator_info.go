package broker

import (
	"context"

	"github.com/bro-n-bro/spacebox/broker/model"
)

func (b *Broker) PublishValidatorInfo(ctx context.Context, info model.ValidatorInfo) error {
	if checkCache(info.ConsensusAddress, info.Height, b.cache.valInfo) {
		return b.marshalAndProduce(ValidatorInfo, info)
	}

	return nil
}
