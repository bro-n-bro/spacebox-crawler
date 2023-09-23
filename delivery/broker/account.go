package broker

import (
	"context"

	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"

	"github.com/bro-n-bro/spacebox/broker/model"
)

func (b *Broker) PublishAccount(ctx context.Context, account model.Account) error {
	data, err := jsoniter.Marshal(account)
	if err != nil {
		return errors.Wrap(err, MsgErrJSONMarshalFail)
	}

	return b.produce(Account, data)
}
