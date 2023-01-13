package broker

import (
	"context"

	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"

	"github.com/hexy-dev/spacebox/broker/model"
)

func (b *Broker) PublishDelegationRewardMessage(ctx context.Context, drm model.DelegationRewardMessage) error {
	data, err := jsoniter.Marshal(drm)
	if err != nil {
		return errors.Wrap(err, MsgErrJSONMarshalFail)
	}

	return b.produce(DelegationRewardMessage, data)
}
