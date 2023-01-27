package broker

import (
	"context"

	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"

	"github.com/bro-n-bro/spacebox/broker/model"
)

func (b *Broker) PublishUnbondingDelegationMessage(ctx context.Context, udm model.UnbondingDelegationMessage) error {
	data, err := jsoniter.Marshal(udm)
	if err != nil {
		return errors.Wrap(err, MsgErrJSONMarshalFail)
	}

	return b.produce(UnbondingDelegationMessage, data)
}
