package broker

import (
	"context"

	"github.com/bro-n-bro/spacebox/broker/model"
)

func (b *Broker) PublishHandleValidatorSignature(ctx context.Context, hvs model.HandleValidatorSignature) error {
	return b.marshalAndProduce(HandleValidatorSignature, hvs)
}
