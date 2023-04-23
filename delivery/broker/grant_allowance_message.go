package broker

import (
	"context"

	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"

	"github.com/bro-n-bro/spacebox/broker/model"
)

func (b *Broker) PublishGrantAllowanceMessage(ctx context.Context, gam model.GrantAllowanceMessage) error {
	data, err := jsoniter.Marshal(gam)
	if err != nil {
		return errors.Wrap(err, MsgErrJSONMarshalFail)
	}

	return b.produce(GrantAllowanceMessage, data)
}
