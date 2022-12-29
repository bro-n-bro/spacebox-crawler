package broker

import (
	"context"

	"github.com/hexy-dev/spacebox/broker/model"
	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"
)

func (b *Broker) PublishValidators(ctx context.Context, vals []model.Validator) error {

	for i := 0; i < len(vals); i++ {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
		data, err := jsoniter.Marshal(vals[i]) // FIXME: maybe user another way to encode data
		if err != nil {
			return errors.Wrap(err, MsgErrJSONMarshalFail)
		}

		if err := b.produce(Validator, data); err != nil {
			return errors.Wrap(err, "produce account fail")
		}
	}
	return nil
}
