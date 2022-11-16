package staking

import (
	"context"
	"encoding/hex"
	"time"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	tmctypes "github.com/tendermint/tendermint/rpc/core/types"
	tmtypes "github.com/tendermint/tendermint/types"

	grpcClient "bro-n-bro-osmosis/client/grpc"
	stakingutils "bro-n-bro-osmosis/modules/staking/utils"
	"bro-n-bro-osmosis/modules/utils"
	"bro-n-bro-osmosis/types"
)

func (m *Module) HandleBlock(ctx context.Context, block *types.Block, vals *tmctypes.ResultValidators) error {
	// Update the validators
	validators, err := stakingutils.UpdateValidators(block.Height, m.client.StakingQueryClient, m.cdc)
	if err != nil {
		return err
	}

	// Get the params
	go updateParams(block.Height, m.client.StakingQueryClient)

	// Update the voting powers
	go updateValidatorVotingPower(block.Height, vals)

	// Update the validators statuses
	go updateValidatorsStatus(block.Height, validators, m.cdc)

	// Updated the double sign evidences
	go updateDoubleSignEvidence(block.Height, block.Evidence.Evidence)

	// Update the staking pool
	go updateStakingPool(block.Height, m.client.StakingQueryClient)

	// Update redelegations and unbonding delegations
	go updateElapsedDelegations(block.Height, block.Timestamp, m.client.StakingQueryClient, m.client.BankQueryClient,
		m.enabledModules)

	return nil
}

// updateParams gets the updated params and stores them inside the database
func updateParams(height int64, stakingClient stakingtypes.QueryClient) {
	//log.Debug().Str("module", "staking").Int64("height", height).
	//	Msg("updating params")
	//
	res, err := stakingClient.Params(
		context.Background(),
		&stakingtypes.QueryParamsRequest{},
		grpcClient.GetHeightRequestHeader(height),
	)
	if err != nil {
		//log.Error().Str("module", "staking").Err(err).
		//	Int64("height", height).
		//	Msg("error while getting params")
		return
	}

	_ = res

	// TODO:
	//err = db.SaveStakingParams(types.NewStakingParams(res.Params, height))
	//if err != nil {
	//log.Error().Str("module", "staking").Err(err).
	//	Int64("height", height).
	//	Msg("error while saving params")
	//return
	//}
}

// updateValidatorsStatus updates all validators' statuses
func updateValidatorsStatus(height int64, validators []stakingtypes.Validator, cdc codec.Codec) {
	//log.Debug().Str("module", "staking").Int64("height", height).
	//	Msg("updating validators statuses")

	statuses, err := stakingutils.GetValidatorsStatuses(height, validators, cdc)
	if err != nil {
		//log.Error().Str("module", "staking").Err(err).
		//	Int64("height", height).
		//	Send()
		return
	}

	_ = statuses
	// TODO:
	//err = db.SaveValidatorsStatuses(statuses)
	//if err != nil {
	//	log.Error().Str("module", "staking").Err(err).
	//		Int64("height", height).
	//		Msg("error while saving validators statuses")
	//}
}

// updateValidatorVotingPower fetches and stores into the database all the current validators' voting powers
func updateValidatorVotingPower(height int64, vals *tmctypes.ResultValidators) {
	//log.Debug().Str("module", "staking").Int64("height", height).
	//	Msg("updating validators voting powers")
	//
	votingPowers := stakingutils.GetValidatorsVotingPowers(height, vals)

	_ = votingPowers
	// TODO:
	//err := db.SaveValidatorsVotingPowers(votingPowers)
	//if err != nil {
	//	log.Error().Str("module", "staking").Err(err).Int64("height", height).
	//		Msg("error while saving validators voting powers")
	//}
}

// updateDoubleSignEvidence updates the double sign evidence of all validators
func updateDoubleSignEvidence(height int64, evidenceList tmtypes.EvidenceList) {
	//log.Debug().Str("module", "staking").Int64("height", height).
	//	Msg("updating double sign evidence")
	//
	for _, ev := range evidenceList {
		dve, ok := ev.(*tmtypes.DuplicateVoteEvidence)
		if !ok {
			continue
		}

		evidence := types.NewDoubleSignEvidence(
			height,
			types.NewDoubleSignVote(
				int(dve.VoteA.Type),
				dve.VoteA.Height,
				dve.VoteA.Round,
				dve.VoteA.BlockID.String(),
				sdk.ConsAddress(dve.VoteA.ValidatorAddress).String(),
				dve.VoteA.ValidatorIndex,
				hex.EncodeToString(dve.VoteA.Signature),
			),
			types.NewDoubleSignVote(
				int(dve.VoteB.Type),
				dve.VoteB.Height,
				dve.VoteB.Round,
				dve.VoteB.BlockID.String(),
				sdk.ConsAddress(dve.VoteB.ValidatorAddress).String(),
				dve.VoteB.ValidatorIndex,
				hex.EncodeToString(dve.VoteB.Signature),
			),
		)

		// TODO:
		_ = evidence
		//err := db.SaveDoubleSignEvidence(evidence)
		//if err != nil {
		//log.Error().Str("module", "staking").Err(err).Int64("height", height).
		//	Msg("error while saving double sign evidence")
		//return
		//}
	}
}

// updateStakingPool reads from the LCD the current staking pool and stores its value inside the database
func updateStakingPool(height int64, stakingClient stakingtypes.QueryClient) {
	//log.Debug().Str("module", "staking").Int64("height", height).
	//	Msg("updating staking pool")

	pool, err := stakingutils.GetStakingPool(height, stakingClient)
	if err != nil {
		//log.Error().Str("module", "staking").Err(err).Int64("height", height).
		//	Msg("error while getting staking pool")
		return
	}

	_ = pool
	// TODO:
	//err = db.SaveStakingPool(pool)
	//if err != nil {
	//	log.Error().Str("module", "staking").Err(err).Int64("height", height).
	//		Msg("error while saving staking pool")
	//	return
	//}
}

// updateElapsedDelegations updates the redelegations and unbonding delegations that have elapsed
func updateElapsedDelegations(
	height int64, timestamp time.Time,
	stakingClient stakingtypes.QueryClient, bankClient banktypes.QueryClient, enabledModules []string,
) {
	//log.Debug().Str("module", "staking").Int64("height", height).
	//	Msg("updating elapsed redelegations and unbonding delegations")
	//
	//deletedRedelegations, err := db.DeleteCompletedRedelegations(timestamp)
	//if err != nil {
	//	log.Error().Str("module", "staking").Err(err).Int64("height", height).
	//		Msg("error while deleting completed redelegations")
	//	return
	//}

	//deletedUnbondingDelegations, err := db.DeleteCompletedUnbondingDelegations(timestamp)
	//if err != nil {
	//	log.Error().Str("module", "staking").Err(err).Int64("height", height).
	//		Msg("error while deleting completed unbonding delegations")
	//	return
	//}

	var delegators = map[string]bool{}

	// Add all the delegators from the redelegations
	//for _, redelegation := range deletedRedelegations {
	//	if _, ok := delegators[redelegation.DelegatorAddress]; !ok {
	//		delegators[redelegation.DelegatorAddress] = true
	//	}
	//}
	//
	//// Add all the delegators from unbonding delegations
	//for _, delegation := range deletedUnbondingDelegations {
	//	if _, ok := delegators[delegation.DelegatorAddress]; !ok {
	//		delegators[delegation.DelegatorAddress] = true
	//	}
	//}

	// Update the delegations and balances of all the delegators
	for delegator := range delegators {
		stakingutils.RefreshDelegations(height, delegator, stakingClient)
		stakingutils.RefreshBalance(delegator, bankClient)

		if utils.ContainAny(enabledModules, "history") {
			// TODO:
			//err := historyutils.UpdateAccountBalanceHistory(delegator)
			//if err != nil {
			//	log.Error().Str("module", "staking").Err(err).Int64("height", height).
			//		Str("account", delegator).
			//		Msg("error while updating account balance history")
			//	return
			//}
		}
	}
}
