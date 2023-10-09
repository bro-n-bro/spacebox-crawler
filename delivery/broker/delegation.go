package broker

import (
	"context"

	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"

	"github.com/bro-n-bro/spacebox/broker/model"
)

func (b *Broker) PublishDelegation(ctx context.Context, d model.Delegation) error {
	d.IsActive = true
	data, err := jsoniter.Marshal(d)
	if err != nil {
		return errors.Wrap(err, MsgErrJSONMarshalFail)
	}

	return b.produce(Delegation, data)
}

func (b *Broker) PublishDisabledDelegation(ctx context.Context, d model.Delegation) error {
	d.IsActive = false
	data, err := jsoniter.Marshal(d)
	if err != nil {
		return errors.Wrap(err, MsgErrJSONMarshalFail)
	}

	return b.produce(Delegation, data)
}
