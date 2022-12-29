package broker

import (
	"context"

	"github.com/hexy-dev/spacebox/broker/model"
	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"
)

func (b *Broker) PublishProposalTallyResult(ctx context.Context, ptr model.ProposalTallyResult) error {

	data, err := jsoniter.Marshal(ptr) // FIXME: maybe user another way to encode data
	if err != nil {
		return errors.Wrap(err, MsgErrJSONMarshalFail)
	}

	if err := b.produce(ProposalTallyResult, data); err != nil {
		return errors.Wrap(err, "produce delegation_reward_message fail")
	}
	return nil
}
