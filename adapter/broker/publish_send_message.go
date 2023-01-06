package broker

import (
	"context"

	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"

	"github.com/hexy-dev/spacebox/broker/model"
)

func (b *Broker) PublishSendMessage(ctx context.Context, sm model.SendMessage) error {

	data, err := jsoniter.Marshal(sm) // FIXME: maybe user another way to encode data
	if err != nil {
		return errors.Wrap(err, MsgErrJSONMarshalFail)
	}

	if err := b.produce(SendMessage, data); err != nil {
		return errors.Wrap(err, "produce send_message fail")
	}
	return nil
}
