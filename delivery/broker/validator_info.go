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
	valInfoCache *lru.Cache[string, int64]
)

func init() {
	var err error
	valInfoCache, err = lru.New[string, int64](100000)
	if err != nil {
		panic(err)
	}
}

func (b *Broker) PublishValidatorInfo(ctx context.Context, info model.ValidatorInfo) error {
	updated := updateCacheValue[string, int64](
		valInfoCache,
		info.ConsensusAddress,
		info.Height, func(curVal int64) bool { return info.Height > curVal })

	if !updated {
		return nil
	}

	data, err := jsoniter.Marshal(info)
	if err != nil {
		return errors.Wrap(err, MsgErrJSONMarshalFail)
	}

	return b.produce(ValidatorInfo, data)
}
