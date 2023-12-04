package broker

import (
	"context"

	"github.com/bro-n-bro/spacebox/broker/model"
)

func (b *Broker) PublishWithdrawValidatorCommissionMessage(
	_ context.Context,
	msg model.WithdrawValidatorCommissionMessage,
) error {

	return b.marshalAndProduce(WithdrawValidatorCommissionMessage, msg)
}
