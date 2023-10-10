package broker

import (
	"context"

	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"

	"github.com/bro-n-bro/spacebox/broker/model"
)

func (b *Broker) PublishCyberlinkMessage(ctx context.Context, msg model.CyberlinkMessage) error {
	data, err := jsoniter.Marshal(msg)
	if err != nil {
		return errors.Wrap(err, MsgErrJSONMarshalFail)
	}

	return b.produce(CyberlinkMessage, data)
}
