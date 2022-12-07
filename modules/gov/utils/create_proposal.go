package utils

import (
	"bro-n-bro-osmosis/internal/rep"
	tb "bro-n-bro-osmosis/pkg/mapper/to_broker"
	"bro-n-bro-osmosis/types"
	"context"

	govtypesv1beta1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
)

// SaveProposals save proposals from genesis file
func SaveProposals(ctx context.Context, sdkProposals govtypesv1beta1.Proposals, broker rep.Broker, mapper tb.ToBroker) error {
	proposals := make([]types.Proposal, len(sdkProposals))
	tallyResults := make([]types.TallyResult, len(sdkProposals))
	deposits := make([]types.ProposalDeposit, len(sdkProposals))

	for index, proposal := range sdkProposals {
		// Since it's not possible to get the proposer, set it to nil
		proposals[index] = types.NewProposal(
			proposal.ProposalId,
			proposal.ProposalRoute(),
			proposal.ProposalType(),
			"",
			proposal.Status.String(),
			proposal.GetContent(),
			proposal.SubmitTime,
			proposal.DepositEndTime,
			proposal.VotingStartTime,
			proposal.VotingEndTime,
		)

		tallyResults[index] = types.NewTallyResult(
			proposal.ProposalId,
			proposal.FinalTallyResult.Yes.Int64(),
			proposal.FinalTallyResult.Abstain.Int64(),
			proposal.FinalTallyResult.No.Int64(),
			proposal.FinalTallyResult.NoWithVeto.Int64(),
			1,
		)

		deposits[index] = types.NewProposalDeposit(
			proposal.ProposalId,
			"",
			proposal.TotalDeposit,
			1,
		)
	}

	for _, pTally := range tallyResults {
		// TODO: test it
		if err := broker.PublishProposalTallyResult(ctx, mapper.MapProposalTallyResult(pTally)); err != nil {
			return err
		}
	}

	// TODO:

	//// Save the proposals
	//err := db.SaveProposals(proposals)
	//if err != nil {
	//	return nil
	//}
	//
	//// Save the deposits
	//err = db.SaveDeposits(deposits)
	//if err != nil {
	//	return nil
	//}
	//
	//// Save the tally results
	//return db.SaveTallyResults(tallyResults)
	return nil
}
