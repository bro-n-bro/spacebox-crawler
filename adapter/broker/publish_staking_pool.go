package broker

import (
	"context"

	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"

	"github.com/hexy-dev/spacebox/broker/model"
)

func (b *Broker) PublishStakingPool(ctx context.Context, sp model.StakingPool) error {
	data, err := jsoniter.Marshal(sp) // FIXME: maybe user another way to encode data
	if err != nil {
		return errors.Wrap(err, MsgErrJSONMarshalFail)
	}

	if err := b.produce(StakingPool, data); err != nil {
		return errors.Wrap(err, "produce stakingTopics pool fail")
	}

	return nil
}
