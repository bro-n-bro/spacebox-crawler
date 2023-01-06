package utils

import (
	"context"
	"fmt"

	"github.com/hexy-dev/spacebox-crawler/types"

	"github.com/hexy-dev/spacebox/broker/model"

	"github.com/cosmos/cosmos-sdk/codec"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	grpcClient "github.com/hexy-dev/spacebox-crawler/client/grpc"
	tb "github.com/hexy-dev/spacebox-crawler/pkg/mapper/to_broker"
	tmctypes "github.com/tendermint/tendermint/rpc/core/types"
	"google.golang.org/grpc/codes"
)

const (
	ErrProposalNotFound = "rpc error: code = %s desc = rpc error: code = %s desc = proposal %d doesn't exist: key not found"
)

type UpdateProposalClients struct {
	govClient     govtypes.QueryClient
	bankClient    banktypes.QueryClient
	stakingClient stakingtypes.QueryClient
}

func NewUpdateProposalClients(gov govtypes.QueryClient, bank banktypes.QueryClient,
	staking stakingtypes.QueryClient) *UpdateProposalClients {
	return &UpdateProposalClients{
		govClient:     gov,
		bankClient:    bank,
		stakingClient: staking,
	}
}

func UpdateProposal(
	ctx context.Context,
	height int64,
	blockVals *tmctypes.ResultValidators,
	id uint64,
	clients *UpdateProposalClients,
	cdc codec.Codec,
	mapper tb.ToBroker,
	broker interface {
		PublishProposal(ctx context.Context, proposal model.Proposal) error
		PublishProposalTallyResult(ctx context.Context, ptr model.ProposalTallyResult) error
	},
) error {

	// Get the proposal
	res, err := clients.govClient.Proposal(ctx, &govtypes.QueryProposalRequest{ProposalId: id})
	if err != nil {
		// Get the error code
		var code string
		_, err = fmt.Sscanf(err.Error(), ErrProposalNotFound, &code, &code, &id)
		if err != nil {
			return err
		}

		if code == codes.NotFound.String() {
			// Handle case when a proposal is deleted from the chain (did not pass deposit period)
			// TODO: delete proposal
			return nil
		}

		return fmt.Errorf("error while getting proposal: %s", err)
	}

	// Unpack the content
	var content govtypes.Content
	err = cdc.UnpackAny(res.Proposal.Content, &content)
	if err != nil {
		return err
	}

	contentBytes, err := types.GetProposalContentBytes(content, cdc)
	if err != nil {
		return err
	}

	proposal := model.NewProposal(
		res.Proposal.ProposalId, content.GetTitle(), content.GetDescription(),
		res.Proposal.ProposalRoute(), res.Proposal.ProposalType(), "", /* FIXME!*/
		res.Proposal.Status.String(), contentBytes,
		res.Proposal.SubmitTime, res.Proposal.DepositEndTime, res.Proposal.VotingStartTime, res.Proposal.VotingEndTime)

	if err = broker.PublishProposal(ctx, proposal); err != nil {
		return fmt.Errorf("error while updating proposal status: %s", err)
	}

	if err = updateProposalTallyResult(ctx, height, res.Proposal, clients.govClient, broker); err != nil {
		return fmt.Errorf("error while updating proposal tally result: %s", err)
	}

	err = updateAccounts(res.Proposal, clients.bankClient)
	if err != nil {
		return fmt.Errorf("error while updating account: %s", err)
	}

	return nil
}

// updateProposalTallyResult updates the tally result associated with the given proposal
func updateProposalTallyResult(
	ctx context.Context,
	height int64,
	proposal govtypes.Proposal,
	govClient govtypes.QueryClient,
	broker interface {
		PublishProposalTallyResult(ctx context.Context, ptr model.ProposalTallyResult) error
	},
) error {

	header := grpcClient.GetHeightRequestHeader(height)
	res, err := govClient.TallyResult(
		context.Background(),
		&govtypes.QueryTallyResultRequest{ProposalId: proposal.ProposalId},
		header,
	)
	if err != nil {
		return err
	}

	tr := model.NewProposalTallyResult(
		proposal.ProposalId,
		height,
		res.Tally.Yes.Int64(),
		res.Tally.Abstain.Int64(),
		res.Tally.No.Int64(),
		res.Tally.NoWithVeto.Int64(),
	)
	// TODO: test it
	err = broker.PublishProposalTallyResult(ctx, tr)
	if err != nil {
		return err
	}
	return nil
}

// updateAccounts updates any account that might be involved in the proposal (eg. fund community recipient)
func updateAccounts(proposal govtypes.Proposal, bankClient banktypes.QueryClient) error {
	// TODO:
	// content, ok := proposal.Content.GetCachedValue().(*distrtypes.CommunityPoolSpendProposal)
	// if ok {
	//	height, err := db.GetLastBlockHeight()
	//	if err != nil {
	//		return err
	//	}
	//
	//	addresses := []string{content.Recipient}
	//
	//	err = authutils.UpdateAccounts(addresses, db)
	//	if err != nil {
	//		return err
	//	}
	//
	//	return bankutils.UpdateBalances(addresses, height, bankClient, db)
	// }
	return nil
}
