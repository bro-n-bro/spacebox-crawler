package broker

import (
	"context"

	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"

	"github.com/hexy-dev/spacebox/broker/model"
)

func (b *Broker) PublishAnnualProvision(ctx context.Context, ap model.AnnualProvision) error {

	data, err := jsoniter.Marshal(ap) // FIXME: maybe user another way to encode data
	if err != nil {
		return errors.Wrap(err, MsgErrJSONMarshalFail)
	}

	if err := b.produce(AnnualProvision, data); err != nil {
		return errors.Wrap(err, "produce block fail")
	}
	return nil
}
