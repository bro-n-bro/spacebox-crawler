package broker

import (
	"context"

	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"

	"github.com/bro-n-bro/spacebox/broker/model"
)

func (b *Broker) PublishProposalVoteMessage(ctx context.Context, pvm model.ProposalVoteMessage) error {
	data, err := jsoniter.Marshal(pvm)
	if err != nil {
		return errors.Wrap(err, MsgErrJSONMarshalFail)
	}

	return b.produce(ProposalVoteMessage, data)
}
