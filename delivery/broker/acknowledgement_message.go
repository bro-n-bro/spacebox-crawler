package broker

import (
	"context"

	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"

	"github.com/bro-n-bro/spacebox/broker/model"
)

func (b *Broker) PublishAcknowledgementMessage(ctx context.Context, msg model.AcknowledgementMessage) error {
	data, err := jsoniter.Marshal(msg)
	if err != nil {
		return errors.Wrap(err, MsgErrJSONMarshalFail)
	}

	return b.produce(AcknowledgementMessage, data)
}
