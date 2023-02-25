package broker

import (
	"context"

	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"

	"github.com/bro-n-bro/spacebox/broker/model"
)

func (b *Broker) PublishValidator(ctx context.Context, val model.Validator) error {
	if !checkCache(val.ConsensusAddress, val.Height, b.cache.validator) {
		return nil
	}

	data, err := jsoniter.Marshal(val)
	if err != nil {
		return errors.Wrap(err, MsgErrJSONMarshalFail)
	}

	return b.produce(Validator, data)
}
