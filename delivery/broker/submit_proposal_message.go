package broker

import (
	"context"

	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"

	"github.com/bro-n-bro/spacebox/broker/model"
)

func (b *Broker) PublishSubmitProposalMessage(_ context.Context, spm model.SubmitProposalMessage) error {
	data, err := jsoniter.Marshal(spm)
	if err != nil {
		return errors.Wrap(err, MsgErrJSONMarshalFail)
	}

	return b.produce(SubmitProposalMessage, data)
}
