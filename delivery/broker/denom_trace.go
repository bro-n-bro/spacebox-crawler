package broker

import (
	"context"

	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"

	"github.com/bro-n-bro/spacebox/broker/model"
)

func (b *Broker) PublishDenomTrace(ctx context.Context, dt model.DenomTrace) error {
	data, err := jsoniter.Marshal(dt)
	if err != nil {
		return errors.Wrap(err, MsgErrJSONMarshalFail)
	}

	return b.produce(DenomTrace, data)
}
