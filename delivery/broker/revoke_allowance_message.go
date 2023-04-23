package broker

import (
	"context"

	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"

	"github.com/bro-n-bro/spacebox/broker/model"
)

func (b *Broker) PublishRevokeAllowanceMessage(ctx context.Context, ram model.RevokeAllowanceMessage) error {
	data, err := jsoniter.Marshal(ram)
	if err != nil {
		return errors.Wrap(err, MsgErrJSONMarshalFail)
	}

	return b.produce(GrantAllowanceMessage, data)
}
