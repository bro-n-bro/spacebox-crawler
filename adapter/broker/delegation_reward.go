package broker

import (
	"context"

	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"

	"github.com/hexy-dev/spacebox/broker/model"
)

func (b *Broker) PublishDelegationReward(ctx context.Context, dr model.DelegationReward) error {
	data, err := jsoniter.Marshal(dr)
	if err != nil {
		return errors.Wrap(err, MsgErrJSONMarshalFail)
	}

	return b.produce(DelegationReward, data)
}
