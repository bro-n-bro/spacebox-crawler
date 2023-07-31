package broker

import (
	"context"

	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"

	"github.com/bro-n-bro/spacebox/broker/model"
)

func (b *Broker) PublishLiquidityPool(_ context.Context, v model.LiquidityPool) error {
	data, err := jsoniter.Marshal(v)
	if err != nil {
		return errors.Wrap(err, MsgErrJSONMarshalFail)
	}

	return b.produce(LiquidityPool, data)
}
