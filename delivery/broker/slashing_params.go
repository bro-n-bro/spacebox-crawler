package broker

import (
	"context"

	"github.com/bro-n-bro/spacebox/broker/model"
)

func (b *Broker) PublishSlashingParams(_ context.Context, params model.SlashingParams) error {
	return b.marshalAndProduce(SlashingParams, params)
}
