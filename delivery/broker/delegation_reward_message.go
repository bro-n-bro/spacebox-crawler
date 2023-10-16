package broker

import (
	"context"

	"github.com/bro-n-bro/spacebox/broker/model"
)

func (b *Broker) PublishDelegationRewardMessage(ctx context.Context, drm model.DelegationRewardMessage) error {
	return b.marshalAndProduce(DelegationRewardMessage, drm)
}
