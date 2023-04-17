package broker

import (
	"context"

	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"

	"github.com/bro-n-bro/spacebox/broker/model"
)

func (b *Broker) PublishGrantMessage(ctx context.Context, gm model.GrantMessage) error {
	data, err := jsoniter.Marshal(gm)
	if err != nil {
		return errors.Wrap(err, MsgErrJSONMarshalFail)
	}

	return b.produce(GrantMessage, data)
}
