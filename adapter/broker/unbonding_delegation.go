package broker

import (
	"context"

	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"

	"github.com/hexy-dev/spacebox/broker/model"
)

func (b *Broker) PublishUnbondingDelegation(ctx context.Context, ud model.UnbondingDelegation) error {
	data, err := jsoniter.Marshal(ud)
	if err != nil {
		return errors.Wrap(err, MsgErrJSONMarshalFail)
	}

	return b.produce(UnbondingDelegation, data)
}
