package broker

import (
	"context"

	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"

	"github.com/hexy-dev/spacebox/broker/model"
)

func (b *Broker) PublishValidatorsInfo(ctx context.Context, infos []model.ValidatorInfo) error {

	for i := 0; i < len(infos); i++ {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
		data, err := jsoniter.Marshal(infos[i]) // FIXME: maybe user another way to encode data
		if err != nil {
			return errors.Wrap(err, MsgErrJSONMarshalFail)
		}

		if err := b.produce(ValidatorInfo, data); err != nil {
			return errors.Wrap(err, "produce account fail")
		}
	}
	return nil
}
