package broker

import (
	"context"

	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"

	"github.com/bro-n-bro/spacebox/broker/model"
)

func (b *Broker) PublishCancelUnbondingDelegationMessage(
	_ context.Context,
	description model.CancelUnbondingDelegationMessage,
) error {

	data, err := jsoniter.Marshal(description)
	if err != nil {
		return errors.Wrap(err, MsgErrJSONMarshalFail)
	}

	return b.produce(CancelUnbondingDelegationMessage, data)
}
