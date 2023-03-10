package broker

import (
	"context"

	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"

	"github.com/bro-n-bro/spacebox/broker/model"
)

func (b *Broker) PublishValidatorInfo(ctx context.Context, info model.ValidatorInfo) error {
	if !checkCache(info.ConsensusAddress, info.Height, b.cache.valInfo) {
		return nil
	}

	data, err := jsoniter.Marshal(info)
	if err != nil {
		return errors.Wrap(err, MsgErrJSONMarshalFail)
	}

	return b.produce(ValidatorInfo, data)
}
