package broker

import (
	"context"

	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"

	"github.com/bro-n-bro/spacebox/broker/model"
)

func (b *Broker) PublishDistributionParams(ctx context.Context, dp model.DistributionParams) error {
	data, err := jsoniter.Marshal(dp)
	if err != nil {
		return errors.Wrap(err, MsgErrJSONMarshalFail)
	}

	return b.produce(DistributionParams, data)
}
