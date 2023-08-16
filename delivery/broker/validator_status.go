package broker

import (
	"context"

	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"

	"github.com/bro-n-bro/spacebox/broker/model"
)

func (b *Broker) PublishValidatorStatus(ctx context.Context, status model.ValidatorStatus) error {
	if !checkCache(status.ConsensusAddress, status.Height, b.cache.valStatus) {
		return nil
	}

	data, err := jsoniter.Marshal(status)
	if err != nil {
		return errors.Wrap(err, MsgErrJSONMarshalFail)
	}

	return b.produce(ValidatorStatus, data)
}
