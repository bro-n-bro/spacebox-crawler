package broker

import (
	"context"

	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"

	"github.com/hexy-dev/spacebox/broker/model"
)

func (b *Broker) PublishProposal(ctx context.Context, proposal model.Proposal) error {

	data, err := jsoniter.Marshal(proposal) // FIXME: maybe user another way to encode data
	if err != nil {
		return errors.Wrap(err, MsgErrJSONMarshalFail)
	}

	if err := b.produce(Proposal, data); err != nil {
		return errors.Wrap(err, "produce proposal fail")
	}
	return nil
}
