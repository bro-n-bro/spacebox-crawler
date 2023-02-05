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
	valCommissionCache *lru.Cache[string, int64]
)

func init() {
	var err error
	valCommissionCache, err = lru.New[string, int64](100000)
	if err != nil {
		panic(err)
	}
}

func (b *Broker) PublishValidatorCommission(ctx context.Context, commission model.ValidatorCommission) error {
	updated := updateCacheValue[string, int64](
		valCommissionCache,
		commission.OperatorAddress,
		commission.Height, func(curVal int64) bool { return commission.Height > curVal })

	if !updated {
		return nil
	}

	data, err := jsoniter.Marshal(commission)
	if err != nil {
		return errors.Wrap(err, MsgErrJSONMarshalFail)
	}

	return b.produce(ValidatorCommission, data)
}
