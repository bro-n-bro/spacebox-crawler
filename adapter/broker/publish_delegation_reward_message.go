package broker

import (
	"context"

	"github.com/hexy-dev/spacebox/broker/model"
	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"
)

func (b *Broker) PublishDelegationRewardMessage(ctx context.Context, drm model.DelegationRewardMessage) error {

	data, err := jsoniter.Marshal(drm) // FIXME: maybe user another way to encode data
	if err != nil {
		return errors.Wrap(err, MsgErrJsonMarshalFail)
	}

	if err := b.produce(DelegationRewardMessage, data); err != nil {
		return errors.Wrap(err, "produce delegation_reward_message fail")
	}
	return nil
}
