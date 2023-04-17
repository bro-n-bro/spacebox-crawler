package broker

import (
	"context"

	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"

	"github.com/bro-n-bro/spacebox/broker/model"
)

func (b *Broker) PublishRevokeMessage(ctx context.Context, rm model.RevokeMessage) error {
	data, err := jsoniter.Marshal(rm)
	if err != nil {
		return errors.Wrap(err, MsgErrJSONMarshalFail)
	}

	return b.produce(RevokeMessage, data)
}
