package broker

import (
	"context"

	"github.com/bro-n-bro/spacebox/broker/model"
)

func (b *Broker) PublishProposalDeposit(ctx context.Context, pvm model.ProposalDeposit) error {
	return b.marshalAndProduce(ProposalDeposit, pvm)
}
