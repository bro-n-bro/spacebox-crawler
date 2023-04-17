package broker

import (
	"context"

	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"

	"github.com/bro-n-bro/spacebox/broker/model"
)

func (b *Broker) PublishCreateValidatorMessage(_ context.Context, cvm model.CreateValidatorMessage) error {
	data, err := jsoniter.Marshal(cvm)
	if err != nil {
		return errors.Wrap(err, MsgErrJSONMarshalFail)
	}

	return b.produce(CreateValidatorMessage, data)
}
