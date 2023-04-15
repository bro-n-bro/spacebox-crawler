package bank

import (
	"context"

	"github.com/bro-n-bro/spacebox/broker/model"
)

type broker interface {
	PublishProposal(ctx context.Context, proposal model.Proposal) error
	PublishGovParams(ctx context.Context, params model.GovParams) error
	PublishProposalDeposit(ctx context.Context, pvm model.ProposalDeposit) error
	PublishProposalDepositMessage(ctx context.Context, pvm model.ProposalDepositMessage) error
	PublishProposalVoteMessage(context.Context, model.ProposalVoteMessage) error
	PublishProposalTallyResult(ctx context.Context, ptr model.ProposalTallyResult) error
	PublishSubmitProposalMessage(ctx context.Context, spm model.SubmitProposalMessage) error
}
