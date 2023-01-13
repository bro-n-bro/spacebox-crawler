package broker

import (
	"context"

	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"

	"github.com/hexy-dev/spacebox/broker/model"
)

func (b *Broker) PublishAccountBalance(ctx context.Context, ab model.AccountBalance) error {
	data, err := jsoniter.Marshal(ab)
	if err != nil {
		return errors.Wrap(err, MsgErrJSONMarshalFail)
	}

	return b.produce(AccountBalance, data)
}
