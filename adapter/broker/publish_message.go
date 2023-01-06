package broker

import (
	"context"

	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"

	"github.com/hexy-dev/spacebox/broker/model"
)

func (b *Broker) PublishMessage(ctx context.Context, message model.Message) error {

	data, err := jsoniter.Marshal(message) // FIXME: maybe user another way to encode data
	if err != nil {
		return errors.Wrap(err, MsgErrJSONMarshalFail)
	}

	if err := b.produce(Message, data); err != nil {
		return errors.Wrap(err, "produce message fail")
	}
	return nil
}
