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
	valStatusCache *lru.Cache[string, int64]
)

func init() {
	var err error
	valStatusCache, err = lru.New[string, int64](100000)
	if err != nil {
		panic(err)
	}
}

func (b *Broker) PublishValidatorStatus(ctx context.Context, status model.ValidatorStatus) error {
	updated := updateCacheValue[string, int64](
		valStatusCache,
		status.ValidatorAddress,
		status.Height, func(curVal int64) bool { return status.Height > curVal })

	if !updated {
		return nil
	}

	data, err := jsoniter.Marshal(status)
	if err != nil {
		return errors.Wrap(err, MsgErrJSONMarshalFail)
	}

	return b.produce(ValidatorStatus, data)
}
