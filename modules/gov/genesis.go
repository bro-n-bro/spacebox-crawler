package bank

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hexy-dev/spacebox-crawler/types"
	"github.com/hexy-dev/spacebox/broker/model"

	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	govtypesv1beta1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	tmtypes "github.com/tendermint/tendermint/types"
)

func (m *Module) HandleGenesis(ctx context.Context, _ *tmtypes.GenesisDoc, appState map[string]json.RawMessage) error {

	// Read the genesis state
	var genState govtypesv1beta1.GenesisState
	err := m.cdc.UnmarshalJSON(appState[govtypes.ModuleName], &genState)
	if err != nil {
		return fmt.Errorf("error while reading gov genesis data: %s", err)
	}

	proposals := genState.Proposals
	if err = m.publishProposals(ctx, proposals); err != nil {
		return fmt.Errorf("error while saving genesis proposal data: %s", err)
	}

	return nil
}

// publishProposals publishes proposals from genesis state to the broker.
func (m *Module) publishProposals(ctx context.Context, proposals govtypesv1beta1.Proposals) error {
	for _, proposal := range proposals {
		// Since it's not possible to get the proposer, set it to nil
		content := proposal.GetContent()
		contentBytes, err := types.GetProposalContentBytes(content, m.cdc)
		if err != nil {
			return err
		}

		// TODO: test it
		if err = m.broker.PublishProposal(ctx,
			model.NewProposal(
				proposal.ProposalId, content.GetTitle(), content.GetDescription(), content.ProposalRoute(),
				content.ProposalType(), "", proposal.Status.String(), contentBytes,
				proposal.SubmitTime, proposal.DepositEndTime, proposal.VotingStartTime, proposal.VotingEndTime),
		); err != nil {
			return err
		}

		// TODO: test it
		// Publish tally results
		if err := m.broker.PublishProposalTallyResult(ctx,
			model.NewProposalTallyResult(
				proposal.ProposalId,
				1,
				proposal.FinalTallyResult.Yes.Int64(),
				proposal.FinalTallyResult.Abstain.Int64(),
				proposal.FinalTallyResult.No.Int64(),
				proposal.FinalTallyResult.NoWithVeto.Int64())); err != nil {

			return err
		}

		// Publish proposal deposits
		if err := m.broker.PublishProposalDeposit(ctx,
			model.NewProposalDeposit(
				proposal.ProposalId,
				1,
				"",
				m.tbM.MapCoins(types.NewCoinsFromCdk(proposal.TotalDeposit)))); err != nil {
			return err
		}

	}

	return nil
}
