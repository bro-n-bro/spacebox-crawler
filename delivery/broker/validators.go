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
	valCache *lru.Cache[string, int64]
)

func init() {
	var err error
	valCache, err = lru.New[string, int64](100000)
	if err != nil {
		panic(err)
	}
}

func (b *Broker) PublishValidator(ctx context.Context, val model.Validator) error {
	updated := updateCacheValue[string, int64](
		valCache,
		val.ConsensusAddress,
		val.Height, func(curVal int64) bool { return val.Height > curVal })

	if !updated {
		return nil
	}

	data, err := jsoniter.Marshal(val)
	if err != nil {
		return errors.Wrap(err, MsgErrJSONMarshalFail)
	}

	return b.produce(Validator, data)
}
