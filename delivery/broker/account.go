package broker

import (
	"context"

	lru "github.com/hashicorp/golang-lru/v2"
	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"

	"github.com/bro-n-bro/spacebox/broker/model"
)

var (
	// TODO: use redis
	accCache *lru.Cache[string, int64]
)

func init() {
	var err error
	accCache, err = lru.New[string, int64](100000)
	if err != nil {
		panic(err)
	}
}

func (b *Broker) PublishAccount(ctx context.Context, account model.Account) error {
	updated := updateCacheValue[string, int64](
		accCache,
		account.Address,
		account.Height, func(curVal int64) bool { return account.Height < curVal })

	if !updated {
		return nil
	}

	data, err := jsoniter.Marshal(account)
	if err != nil {
		return errors.Wrap(err, MsgErrJSONMarshalFail)
	}

	return b.produce(Account, data)
}
