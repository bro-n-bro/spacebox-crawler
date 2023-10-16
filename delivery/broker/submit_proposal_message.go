package broker

import (
	"context"

	"github.com/bro-n-bro/spacebox/broker/model"
)

func (b *Broker) PublishSubmitProposalMessage(_ context.Context, spm model.SubmitProposalMessage) error {
	return b.marshalAndProduce(SubmitProposalMessage, spm)
}
