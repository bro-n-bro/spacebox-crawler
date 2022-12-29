package broker

import (
	"context"

	"github.com/pkg/errors"

	"github.com/hexy-dev/spacebox/broker/model"

	jsoniter "github.com/json-iterator/go"
)

func (b *Broker) PublishRedelegation(ctx context.Context, r model.Redelegation) error {

	data, err := jsoniter.Marshal(r) // FIXME: maybe user another way to encode data
	if err != nil {
		return errors.Wrap(err, MsgErrJSONMarshalFail)
	}

	if err := b.produce(Redelegation, data); err != nil {
		return err
	}
	return nil
}
