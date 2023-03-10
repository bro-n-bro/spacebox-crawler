package broker

import (
	"context"

	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"

	"github.com/bro-n-bro/spacebox/broker/model"
)

func (b *Broker) PublishAnnualProvision(ctx context.Context, ap model.AnnualProvision) error {
	data, err := jsoniter.Marshal(ap)
	if err != nil {
		return errors.Wrap(err, MsgErrJSONMarshalFail)
	}

	return b.produce(AnnualProvision, data)
}
