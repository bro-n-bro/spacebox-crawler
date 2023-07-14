package broker

import (
	"context"

	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"

	"github.com/bro-n-bro/spacebox/broker/model"
)

func (b *Broker) PublishReceivePacketMessage(ctx context.Context, r model.RecvPacketMessage) error {
	data, err := jsoniter.Marshal(r)
	if err != nil {
		return errors.Wrap(err, MsgErrJSONMarshalFail)
	}

	return b.produce(ReceivePacketMessage, data)
}
