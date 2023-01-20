package broker

import (
	"context"

	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"

	"github.com/bro-n-bro/spacebox/broker/model"
)

func (b *Broker) PublishProposalDepositMessage(ctx context.Context, pvm model.ProposalDepositMessage) error {
	data, err := jsoniter.Marshal(pvm)
	if err != nil {
		return errors.Wrap(err, MsgErrJSONMarshalFail)
	}

	return b.produce(ProposalDepositMessage, data)
}
