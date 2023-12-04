package broker

import (
	"context"

	"github.com/bro-n-bro/spacebox/broker/model"
)

func (b *Broker) PublishProposalTallyResult(ctx context.Context, ptr model.ProposalTallyResult) error {
	return b.marshalAndProduce(ProposalTallyResult, ptr)
}
