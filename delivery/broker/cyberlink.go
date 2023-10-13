package broker

import (
	"context"

	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"

	"github.com/bro-n-bro/spacebox/broker/model"
)

func (b *Broker) PublishCyberLink(ctx context.Context, msg model.CyberLink) error {
	data, err := jsoniter.Marshal(msg)
	if err != nil {
		return errors.Wrap(err, MsgErrJSONMarshalFail)
	}

	return b.produce(CyberLink, data)
}

func (b *Broker) PublishCyberLinkMessage(ctx context.Context, msg model.CyberLinkMessage) error {
	data, err := jsoniter.Marshal(msg)
	if err != nil {
		return errors.Wrap(err, MsgErrJSONMarshalFail)
	}

	return b.produce(CyberLinkMessage, data)
}
