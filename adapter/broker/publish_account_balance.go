package broker

import (
	"context"

	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"

	"github.com/hexy-dev/spacebox/broker/model"
)

func (b *Broker) PublishAccountBalance(ctx context.Context, ab model.AccountBalance) error {

	data, err := jsoniter.Marshal(ab) // FIXME: maybe user another way to encode data
	if err != nil {
		return errors.Wrap(err, MsgErrJSONMarshalFail)
	}
	if err := b.produce(AccountBalance, data); err != nil {
		return errors.Wrap(err, "produce supply fail")
	}
	return nil
}
