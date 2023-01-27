package broker

import (
	"context"

	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"

	"github.com/bro-n-bro/spacebox/broker/model"
)

func (b *Broker) PublishValidatorDescription(ctx context.Context, description model.ValidatorDescription) error {
	data, err := jsoniter.Marshal(description)
	if err != nil {
		return errors.Wrap(err, MsgErrJSONMarshalFail)
	}

	return b.produce(ValidatorDescription, data)
}
