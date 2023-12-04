package broker

import (
	"context"

	"github.com/bro-n-bro/spacebox/broker/model"
)

func (b *Broker) PublishValidatorVotingPower(_ context.Context, vp model.ValidatorVotingPower) error {
	return b.marshalAndProduce(ValidatorVotingPower, vp)
}
