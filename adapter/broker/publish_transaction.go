package broker

import (
	"context"

	"github.com/hexy-dev/spacebox/broker/model"
	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"
)

func (b *Broker) PublishTransaction(ctx context.Context, tx model.Transaction) error {

	data, err := jsoniter.Marshal(tx) // FIXME: maybe user another way to encode data
	if err != nil {
		return errors.Wrap(err, MsgErrJSONMarshalFail)
	}

	if err := b.produce(Transaction, data); err != nil {
		return errors.Wrap(err, "produce transaction fail")
	}
	return nil
}
