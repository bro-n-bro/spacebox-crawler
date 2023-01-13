package broker

import (
	"context"

	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"

	"github.com/hexy-dev/spacebox/broker/model"
)

func (b *Broker) PublishDelegationMessage(ctx context.Context, dm model.DelegationMessage) error {
	data, err := jsoniter.Marshal(dm)
	if err != nil {
		return errors.Wrap(err, MsgErrJSONMarshalFail)
	}

	return b.produce(DelegationMessage, data)
}
