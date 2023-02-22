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
	valDescriptionCache *lru.Cache[string, int64]
)

func init() {
	var err error
	valDescriptionCache, err = lru.New[string, int64](100000)
	if err != nil {
		panic(err)
	}
}

func (b *Broker) PublishValidatorDescription(ctx context.Context, description model.ValidatorDescription) error {
	updated := updateCacheValue[string, int64](
		valDescriptionCache,
		description.OperatorAddress,
		description.Height, func(curVal int64) bool { return description.Height > curVal })

	if !updated {
		return nil
	}

	data, err := jsoniter.Marshal(description)
	if err != nil {
		return errors.Wrap(err, MsgErrJSONMarshalFail)
	}

	return b.produce(ValidatorDescription, data)
}
