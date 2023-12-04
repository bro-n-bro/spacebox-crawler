package broker

import (
	"context"

	"github.com/bro-n-bro/spacebox/broker/model"
)

func (b *Broker) PublishProposalVoteMessage(ctx context.Context, pvm model.ProposalVoteMessage) error {
	return b.marshalAndProduce(ProposalVoteMessage, pvm)
}
