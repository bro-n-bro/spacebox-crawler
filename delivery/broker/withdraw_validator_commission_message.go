package broker

import (
	"context"

	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"

	"github.com/bro-n-bro/spacebox/broker/model"
)

func (b *Broker) PublishWithdrawValidatorCommissionMessage(
	_ context.Context,
	wvcm model.WithdrawValidatorCommissionMessage,
) error {

	data, err := jsoniter.Marshal(wvcm)
	if err != nil {
		return errors.Wrap(err, MsgErrJSONMarshalFail)
	}

	return b.produce(WithdrawValidatorCommissionMessage, data)
}
