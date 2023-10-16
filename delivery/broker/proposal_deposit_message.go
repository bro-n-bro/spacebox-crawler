package broker

import (
	"context"

	"github.com/bro-n-bro/spacebox/broker/model"
)

func (b *Broker) PublishProposalDepositMessage(ctx context.Context, pvm model.ProposalDepositMessage) error {
	return b.marshalAndProduce(ProposalDepositMessage, pvm)
}
