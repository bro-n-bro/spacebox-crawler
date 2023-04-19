package broker

import (
	"context"

	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"

	"github.com/bro-n-bro/spacebox/broker/model"
)

func (b *Broker) PublishDistributionCommission(_ context.Context, commission model.DistributionCommission) error {
	data, err := jsoniter.Marshal(commission)
	if err != nil {
		return errors.Wrap(err, MsgErrJSONMarshalFail)
	}

	return b.produce(DistributionCommission, data)
}
