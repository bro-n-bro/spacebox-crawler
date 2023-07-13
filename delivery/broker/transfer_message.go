package broker

import (
	"context"

	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"

	"github.com/bro-n-bro/spacebox/broker/model"
)

func (b *Broker) PublishTransferMessage(ctx context.Context, tm model.TransferMessage) error {
	data, err := jsoniter.Marshal(tm)
	if err != nil {
		return errors.Wrap(err, MsgErrJSONMarshalFail)
	}

	return b.produce(TransferMessage, data)
}
