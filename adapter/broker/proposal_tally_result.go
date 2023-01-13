package broker

import (
	"context"

	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"

	"github.com/hexy-dev/spacebox/broker/model"
)

func (b *Broker) PublishProposalTallyResult(ctx context.Context, ptr model.ProposalTallyResult) error {
	data, err := jsoniter.Marshal(ptr)
	if err != nil {
		return errors.Wrap(err, MsgErrJSONMarshalFail)
	}

	return b.produce(ProposalTallyResult, data)
}
