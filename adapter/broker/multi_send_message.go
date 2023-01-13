package broker

import (
	"context"

	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"

	"github.com/hexy-dev/spacebox/broker/model"
)

func (b *Broker) PublishMultiSendMessage(ctx context.Context, msm model.MultiSendMessage) error {
	data, err := jsoniter.Marshal(msm)
	if err != nil {
		return errors.Wrap(err, MsgErrJSONMarshalFail)
	}

	return b.produce(MultiSendMessage, data)
}
