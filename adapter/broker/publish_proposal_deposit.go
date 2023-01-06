package broker

import (
	"context"

	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"

	"github.com/hexy-dev/spacebox/broker/model"
)

func (b *Broker) PublishProposalDeposit(ctx context.Context, pvm model.ProposalDeposit) error {

	data, err := jsoniter.Marshal(pvm) // FIXME: maybe user another way to encode data
	if err != nil {
		return errors.Wrap(err, MsgErrJSONMarshalFail)
	}

	if err := b.produce(ProposalDeposit, data); err != nil {
		return errors.Wrap(err, "produce proposal_deposit fail")
	}
	return nil
}
