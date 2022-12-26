package broker

import (
	"context"

	"github.com/hexy-dev/spacebox/broker/model"
	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"
)

func (b *Broker) PublishValidatorsStatuses(ctx context.Context, statuses []model.ValidatorStatus) error {

	for i := 0; i < len(statuses); i++ {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
		data, err := jsoniter.Marshal(statuses[i]) // FIXME: maybe user another way to encode data
		if err != nil {
			return errors.Wrap(err, MsgErrJsonMarshalFail)
		}

		if err := b.produce(ValidatorStatus, data); err != nil {
			return errors.Wrap(err, "produce account fail")
		}
	}
	return nil
}