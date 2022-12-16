package utils

import (
	grpcClient "bro-n-bro-osmosis/client/grpc"
	"bro-n-bro-osmosis/internal/rep"
	stakingutils "bro-n-bro-osmosis/modules/staking/utils"
	tb "bro-n-bro-osmosis/pkg/mapper/to_broker"
	"bro-n-bro-osmosis/types"
	"context"
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
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
	ctx context.Context, height int64, blockVals *tmctypes.ResultValidators, id uint64, clients *UpdateProposalClients,
	cdc codec.Codec, broker rep.Broker, mapper tb.ToBroker,
) error {
	// Get the proposal
	res, err := clients.govClient.Proposal(ctx, &govtypes.QueryProposalRequest{ProposalId: id})
	if err != nil {
		// Get the error code
		var code string
		_, err := fmt.Sscanf(err.Error(), ErrProposalNotFound, &code, &code, &id)
		if err != nil {
			return err
		}

		if code == codes.NotFound.String() {
			// Handle case when a proposal is deleted from the chain (did not pass deposit period)
			return updateDeletedProposalStatus(id)
		}

		return fmt.Errorf("error while getting proposal: %s", err)
	}

	err = updateProposalStatus(res.Proposal)
	if err != nil {
		return fmt.Errorf("error while updating proposal status: %s", err)
	}

	err = updateProposalTallyResult(ctx, height, res.Proposal, clients.govClient, broker, mapper)
	if err != nil {
		return fmt.Errorf("error while updating proposal tally result: %s", err)
	}

	err = updateAccounts(res.Proposal, clients.bankClient)
	if err != nil {
		return fmt.Errorf("error while updating account: %s", err)
	}

	err = updateProposalStakingPoolSnapshot(height, id, clients.stakingClient)
	if err != nil {
		return fmt.Errorf("error while updating proposal staking pool snapshot: %s", err)
	}

	err = updateProposalValidatorStatusesSnapshot(height, id, blockVals, clients.stakingClient, cdc)
	if err != nil {
		return fmt.Errorf("error while updating proposal validator statuses snapshot: %s", err)
	}

	return nil
}

// updateDeletedProposalStatus updates the proposal having the given id by setting its status
// to the one that represents a deleted proposal
func updateDeletedProposalStatus(id uint64) error {
	//stored, err := db.GetProposal(id)
	//if err != nil {
	//	return err
	//}
	//
	// TODO:
	return nil
	//return db.UpdateProposal(
	//	types.NewProposalUpdate(
	//		stored.ProposalID,
	//		types.ProposalStatusInvalid,
	//		stored.VotingStartTime,
	//		stored.VotingEndTime,
	//	),
	//)
}

// updateProposalStatus updates the given proposal status
func updateProposalStatus(proposal govtypes.Proposal) error {
	// TODO:
	return nil
	//return db.UpdateProposal(
	//	types.NewProposalUpdate(
	//		proposal.ProposalId,
	//		proposal.Status.String(),
	//		proposal.VotingStartTime,
	//		proposal.VotingEndTime,
	//	),
	//)
}

// updateProposalTallyResult updates the tally result associated with the given proposal
func updateProposalTallyResult(ctx context.Context, height int64, proposal govtypes.Proposal, govClient govtypes.QueryClient,
	broker rep.Broker, mapper tb.ToBroker) error {

	//height, err := db.GetLastBlockHeight()
	//if err != nil {
	//	return err
	//}

	header := grpcClient.GetHeightRequestHeader(height)
	res, err := govClient.TallyResult(
		context.Background(),
		&govtypes.QueryTallyResultRequest{ProposalId: proposal.ProposalId},
		header,
	)
	if err != nil {
		return err
	}

	tr := types.NewTallyResult(
		proposal.ProposalId,
		res.Tally.Yes.Int64(),
		res.Tally.Abstain.Int64(),
		res.Tally.No.Int64(),
		res.Tally.NoWithVeto.Int64(),
		height,
	)

	// TODO: test it
	err = broker.PublishProposalTallyResult(ctx, mapper.MapProposalTallyResult(tr))
	if err != nil {
		return err
	}
	return nil

	//return db.SaveTallyResults([]types.TallyResult{
	//	types.NewTallyResult(
	//		proposal.ProposalId,
	//		res.Tally.Yes.Int64(),
	//		res.Tally.Abstain.Int64(),
	//		res.Tally.No.Int64(),
	//		res.Tally.NoWithVeto.Int64(),
	//		height,
	//	),
	//})
}

// updateAccounts updates any account that might be involved in the proposal (eg. fund community recipient)
func updateAccounts(proposal govtypes.Proposal, bankClient banktypes.QueryClient) error {
	// TODO:
	//content, ok := proposal.Content.GetCachedValue().(*distrtypes.CommunityPoolSpendProposal)
	//if ok {
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
	//}
	return nil
}

// updateProposalStakingPoolSnapshot updates the staking pool snapshot associated with the gov
// proposal having the provided id
func updateProposalStakingPoolSnapshot(
	height int64, proposalID uint64, stakingClient stakingtypes.QueryClient,
) error {
	pool, err := stakingutils.GetStakingPool(height, stakingClient)
	if err != nil {
		return fmt.Errorf("error while getting staking pool: %s", err)
	}

	// TODO:
	_ = pool
	return nil
	//return db.SaveProposalStakingPoolSnapshot(
	//	types.NewProposalStakingPoolSnapshot(proposalID, pool),
	//)
}

// updateProposalValidatorStatusesSnapshot updates the snapshots of the various validators for
// the proposal having the given id
func updateProposalValidatorStatusesSnapshot(
	height int64, proposalID uint64, blockVals *tmctypes.ResultValidators, stakingClient stakingtypes.QueryClient,
	cdc codec.Codec) error {
	validators, _, err := stakingutils.GetValidatorsWithStatus(height, stakingtypes.Bonded.String(), stakingClient, cdc)
	if err != nil {
		return err
	}

	votingPowers := stakingutils.GetValidatorsVotingPowers(height, blockVals)

	statuses, _, err := stakingutils.GetValidatorsStatuses(height, validators, cdc)
	if err != nil {
		return err
	}

	var snapshots = make([]types.ProposalValidatorStatusSnapshot, len(validators))
	for index, validator := range validators {
		consAddr, err := validator.GetConsAddr()
		if err != nil {
			return err
		}

		status, err := findStatus(consAddr.String(), statuses)
		if err != nil {
			return err
		}

		votingPower, err := findVotingPower(consAddr.String(), votingPowers)
		if err != nil {
			return err
		}

		snapshots[index] = types.NewProposalValidatorStatusSnapshot(
			proposalID,
			consAddr.String(),
			status.Status,
			status.Jailed,
			votingPower.VotingPower,
			height,
		)
	}

	// TODO:
	return nil
	//return db.SaveProposalValidatorsStatusesSnapshots(snapshots)
}

func findVotingPower(consAddr string, powers []types.ValidatorVotingPower) (types.ValidatorVotingPower, error) {
	for _, votingPower := range powers {
		if votingPower.ConsensusAddress == consAddr {
			return votingPower, nil
		}
	}
	return types.ValidatorVotingPower{}, fmt.Errorf("voting power not found for validator with consensus address %s", consAddr)
}

func findStatus(consAddr string, statuses []types.ValidatorStatus) (types.ValidatorStatus, error) {
	for _, status := range statuses {
		if status.ConsensusAddress == consAddr {
			return status, nil
		}
	}
	return types.ValidatorStatus{}, fmt.Errorf("cannot find status for validator with consensus address %s", consAddr)
}
