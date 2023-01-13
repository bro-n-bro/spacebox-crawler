package broker

import (
	"context"

	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"

	"github.com/hexy-dev/spacebox/broker/model"
)

func (b *Broker) PublishBlock(ctx context.Context, block model.Block) error {
	data, err := jsoniter.Marshal(block)
	if err != nil {
		return errors.Wrap(err, MsgErrJSONMarshalFail)
	}

	return b.produce(Block, data)
}
