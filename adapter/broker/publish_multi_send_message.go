package broker

import (
	"context"

	"github.com/hexy-dev/spacebox/broker/model"
	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"
)

func (b *Broker) PublishMultiSendMessage(ctx context.Context, msm model.MultiSendMessage) error {

	data, err := jsoniter.Marshal(msm) // FIXME: maybe user another way to encode data
	if err != nil {
		return errors.Wrap(err, MsgErrJSONMarshalFail)
	}

	if err := b.produce(MultiSendMessage, data); err != nil {
		return errors.Wrap(err, "produce send_message fail")
	}
	return nil
}
