package broker

import (
	"context"

	"github.com/bro-n-bro/spacebox/broker/model"
)

func (b *Broker) PublishBandwidthParams(ctx context.Context, msg model.BandwidthParams) error {
	return b.marshalAndProduce(BandwidthParams, msg)
}
