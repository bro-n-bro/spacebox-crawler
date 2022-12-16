package broker

import (
	"context"

	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"

	"github.com/hexy-dev/spacebox/broker/model"
)

func (b *Broker) PublishBlock(ctx context.Context, block model.Block) error {

	data, err := jsoniter.Marshal(block) // FIXME: maybe user another way to encode data
	if err != nil {
		return errors.Wrap(err, MsgErrJsonMarshalFail)
	}

	if err := b.produce(BlockTopic, data); err != nil {
		return errors.Wrap(err, "produce block fail")
	}

	return nil
}
