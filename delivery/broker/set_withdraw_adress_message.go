package broker

import (
	"context"

	"github.com/bro-n-bro/spacebox/broker/model"
)

func (b *Broker) PublishSetWithdrawAddressMessage(_ context.Context, swm model.SetWithdrawAddressMessage) error {
	return b.marshalAndProduce(SetWithdrawAddressMessage, swm)
}
