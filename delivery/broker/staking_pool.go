package broker

import (
	"context"

	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"

	"github.com/bro-n-bro/spacebox/broker/model"
)

func (b *Broker) PublishStakingPool(ctx context.Context, sp model.StakingPool) error {
	data, err := jsoniter.Marshal(sp)
	if err != nil {
		return errors.Wrap(err, MsgErrJSONMarshalFail)
	}

	return b.produce(StakingPool, data)
}
