package broker

import (
	"context"

	"github.com/hexy-dev/spacebox/broker/model"
	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"
)

func (b *Broker) PublishSupply(ctx context.Context, supply model.Supply) error {

	data, err := jsoniter.Marshal(supply) // FIXME: maybe user another way to encode data
	if err != nil {
		return errors.Wrap(err, MsgErrJsonMarshalFail)
	}

	if err := b.produce(Supply, data); err != nil {
		return errors.Wrap(err, "produce supply fail")
	}
	return nil
}
