package broker

import (
	"context"

	"github.com/bro-n-bro/spacebox/broker/model"
)

func (b *Broker) PublishVoteWeightedMessage(_ context.Context, vwm model.VoteWeightedMessage) error {
	return b.marshalAndProduce(VoteWeightedMessage, vwm)
}
