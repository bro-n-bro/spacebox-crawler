package broker

import (
	"context"

	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"

	"github.com/bro-n-bro/spacebox/broker/model"
)

func (b *Broker) PublishAccount(ctx context.Context, account model.Account) error {
	if !b.needProduceAccount(account.Address, account.Height) {
		return nil
	}

	data, err := jsoniter.Marshal(account)
	if err != nil {
		return errors.Wrap(err, MsgErrJSONMarshalFail)
	}

	return b.produce(Account, data)
}

func (b *Broker) needProduceAccount(address string, height int64) bool {
	if b.cache.account != nil && !b.cache.account.UpdateCacheValue(address, height) {
		return false
	}
	return true
}
