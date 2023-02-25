package broker

import (
	"context"

	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"

	"github.com/bro-n-bro/spacebox/broker/model"
)

func (b *Broker) PublishValidatorCommission(ctx context.Context, commission model.ValidatorCommission) error {
	if !checkCache(commission.OperatorAddress, commission.Height, b.cache.valCommission) {
		return nil
	}

	data, err := jsoniter.Marshal(commission)
	if err != nil {
		return errors.Wrap(err, MsgErrJSONMarshalFail)
	}

	return b.produce(ValidatorCommission, data)
}
