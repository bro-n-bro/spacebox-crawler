package broker

import (
	"context"

	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"

	"github.com/bro-n-bro/spacebox/broker/model"
)

func (b *Broker) PublishHandleValidatorSignature(ctx context.Context, hvs model.HandleValidatorSignature) error {
	data, err := jsoniter.Marshal(hvs)
	if err != nil {
		return errors.Wrap(err, MsgErrJSONMarshalFail)
	}

	return b.produce(HandleValidatorSignature, data)
}
