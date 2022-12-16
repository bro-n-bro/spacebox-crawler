package broker

import (
	"context"

	"github.com/pkg/errors"

	"github.com/hexy-dev/spacebox/broker/model"

	jsoniter "github.com/json-iterator/go"
)

func (b *Broker) PublishCommunityPool(ctx context.Context, cp model.CommunityPool) error {

	data, err := jsoniter.Marshal(cp) // FIXME: maybe user another way to encode data
	if err != nil {
		return errors.Wrap(err, MsgErrJsonMarshalFail)
	}

	if err := b.produce(CommunityPool, data); err != nil {
		return errors.Wrap(err, "produce supply fail")
	}
	return nil
}
