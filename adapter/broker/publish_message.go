package broker

import (
	"context"

	"github.com/hexy-dev/spacebox/broker/model"
	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"
)

func (b *Broker) PublishMessage(ctx context.Context, message model.Message) error {

	data, err := jsoniter.Marshal(message) // FIXME: maybe user another way to encode data
	if err != nil {
		return errors.Wrap(err, MsgErrJsonMarshalFail)
	}

	if err := b.produce(MessageTopic, data); err != nil {
		return errors.Wrap(err, "produce message fail")
	}
	return nil
}
