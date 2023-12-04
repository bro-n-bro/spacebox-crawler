package broker

import (
	"context"

	"github.com/bro-n-bro/spacebox/broker/model"
)

func (b *Broker) PublishProposal(ctx context.Context, proposal model.Proposal) error {
	return b.marshalAndProduce(Proposal, proposal)
}
