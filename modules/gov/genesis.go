package bank

import (
	"context"
	"encoding/json"
	"fmt"

	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	govtypesv1beta1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	tmtypes "github.com/tendermint/tendermint/types"

	"github.com/bro-n-bro/spacebox-crawler/types"
	"github.com/bro-n-bro/spacebox/broker/model"
)

func (m *Module) HandleGenesis(ctx context.Context, _ *tmtypes.GenesisDoc, appState map[string]json.RawMessage) error {
	// Read the genesis state
	var genState govtypesv1beta1.GenesisState
	if err := m.cdc.UnmarshalJSON(appState[govtypes.ModuleName], &genState); err != nil {
		return fmt.Errorf("error while reading gov genesis data: %w", err)
	}

	if err := m.publishProposals(ctx, genState.Proposals); err != nil {
		return fmt.Errorf("error while saving genesis proposal data: %w", err)
	}

	return nil
}

// publishProposals publishes proposals from genesis state to the broker.
func (m *Module) publishProposals(ctx context.Context, proposals govtypesv1beta1.Proposals) error {
	for _, proposal := range proposals {
		// Since it's not possible to get the proposer, set it to nil
		content := proposal.GetContent()

		contentBytes, err := getProposalContentBytes(content, m.cdc)
		if err != nil {
			return err
		}

		// TODO: test it
		if err = m.broker.PublishProposal(ctx, model.Proposal{
			ID:              proposal.ProposalId,
			Title:           content.GetTitle(),
			Description:     content.GetDescription(),
			ProposalRoute:   content.ProposalRoute(),
			ProposalType:    content.ProposalType(),
			ProposerAddress: "",
			Status:          proposal.Status.String(),
			Content:         contentBytes,
			SubmitTime:      proposal.SubmitTime,
			DepositEndTime:  proposal.DepositEndTime,
			VotingStartTime: proposal.VotingStartTime,
			VotingEndTime:   proposal.VotingEndTime,
		}); err != nil {
			return err
		}

		// TODO: test it
		// Publish tally results
		if err := m.broker.PublishProposalTallyResult(ctx, model.ProposalTallyResult{
			ProposalID: proposal.ProposalId,
			Height:     1,
			Yes:        proposal.FinalTallyResult.Yes.Int64(),
			Abstain:    proposal.FinalTallyResult.Abstain.Int64(),
			No:         proposal.FinalTallyResult.No.Int64(),
			NoWithVeto: proposal.FinalTallyResult.NoWithVeto.Int64(),
		}); err != nil {
			return err
		}

		// Publish proposal deposits
		if err := m.broker.PublishProposalDeposit(ctx, model.ProposalDeposit{
			ProposalID:       proposal.ProposalId,
			Height:           1,
			DepositorAddress: "",
			Coins:            m.tbM.MapCoins(types.NewCoinsFromCdk(proposal.TotalDeposit)),
		}); err != nil {
			return err
		}
	}

	return nil
}
