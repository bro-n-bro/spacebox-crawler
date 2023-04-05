package broker

import (
	"context"

	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"

	"github.com/bro-n-bro/spacebox/broker/model"
)

func (b *Broker) PublishSetWithdrawAddressMessage(_ context.Context, swm model.SetWithdrawAddressMessage) error {
	data, err := jsoniter.Marshal(swm)
	if err != nil {
		return errors.Wrap(err, MsgErrJSONMarshalFail)
	}

	return b.produce(SetWithdrawAddressMessage, data)
}
