package staking

import (
	"context"
	"encoding/hex"
	"time"

	"cosmossdk.io/errors"

	"golang.org/x/sync/errgroup"

	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	tmctypes "github.com/tendermint/tendermint/rpc/core/types"
	tmtypes "github.com/tendermint/tendermint/types"

	grpcClient "github.com/hexy-dev/spacebox-crawler/client/grpc"
	stakingutils "github.com/hexy-dev/spacebox-crawler/modules/staking/utils"
	"github.com/hexy-dev/spacebox-crawler/modules/utils"
	"github.com/hexy-dev/spacebox-crawler/types"
)

func (m *Module) HandleBlock(ctx context.Context, block *types.Block, vals *tmctypes.ResultValidators) error {
	// Update the validators
	validators, err := stakingutils.UpdateValidators(ctx, block.Height, m.client.StakingQueryClient, m.cdc, m.broker, m.tbM)
	if err != nil {
		return err
	}

	g, _ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		// Update the params
		return m.updateParams(_ctx, block.Height)
	})

	// Update the voting powers
	// go updateValidatorVotingPower(block.Height, vals)

	g.Go(func() error {
		// Update the validators statuses
		return m.updateValidatorsStatus(_ctx, block.Height, validators)
	})

	// Updated the double sign evidences
	// go updateDoubleSignEvidence(block.Height, block.Evidence.Evidence)

	g.Go(func() error {
		// Update the staking pool
		return m.updateStakingPool(_ctx, block.Height, m.client.StakingQueryClient)
	})

	g.Go(func() error {
		// Update redelegations and unbonding delegations
		// TODO
		return m.updateElapsedDelegations(_ctx, block.Height, block.Timestamp, m.enabledModules)
	})

	if err := g.Wait(); err != nil {
		return err
	}

	return nil
}

// updateParams gets the updated params and stores them inside the database
func (m *Module) updateParams(ctx context.Context, height int64) error {
	res, err := m.client.StakingQueryClient.Params(
		context.Background(),
		&stakingtypes.QueryParamsRequest{},
		grpcClient.GetHeightRequestHeader(height),
	)
	if err != nil {
		return err
	}

	// TODO: test it
	// TODO: maybe check diff from mongo in my side?
	err = m.broker.PublishStakingParams(ctx, m.tbM.MapStakingParams(types.NewStakingParams(res.Params, height)))
	if err != nil {
		return err
	}

	// TODO:
	// err = db.SaveStakingParams(types.NewStakingParams(res.Params, height))
	// if err != nil {
	// log.Error().Str("module", "staking").Err(err).
	//	Int64("height", height).
	//	Msg("error while saving params")
	// return
	//}

	return nil
}

// updateValidatorsStatus updates all validators' statuses
func (m *Module) updateValidatorsStatus(ctx context.Context, height int64, stakingValidators []stakingtypes.Validator) error {
	statuses, validators, err := stakingutils.GetValidatorsStatuses(height, stakingValidators, m.cdc)
	if err != nil {
		return err
	}

	// TODO: save to mongo?
	// TODO: test it
	if err = m.broker.PublishValidators(ctx, m.tbM.MapValidators(validators)); err != nil {
		return err
	}

	// TODO: test it
	if err = m.broker.PublishValidatorsStatuses(ctx, m.tbM.MapValidatorsStatuses(statuses)); err != nil {
		return err
	}

	// if err != nil {
	//	log.Error().Str("module", "staking").Err(err).
	//		Int64("height", height).
	//		Msg("error while saving stakingValidators statuses")
	// }
	return nil
}

// nolint:deadcode,unused
// updateValidatorVotingPower fetches and stores into the database all the current validators' voting powers
func updateValidatorVotingPower(height int64, vals *tmctypes.ResultValidators) {
	votingPowers := stakingutils.GetValidatorsVotingPowers(height, vals)

	_ = votingPowers
	// TODO:
	// err := db.SaveValidatorsVotingPowers(votingPowers)
	// if err != nil {
	//	log.Error().Str("module", "staking").Err(err).Int64("height", height).
	//		Msg("error while saving validators voting powers")
	// }
}

// nolint:deadcode,unused
// updateDoubleSignEvidence updates the double sign evidence of all validators
func updateDoubleSignEvidence(height int64, evidenceList tmtypes.EvidenceList) {
	// log.Debug().Str("module", "staking").Int64("height", height).
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
		// err := db.SaveDoubleSignEvidence(evidence)
		// if err != nil {
		// log.Error().Str("module", "staking").Err(err).Int64("height", height).
		//	Msg("error while saving double sign evidence")
		// return
		// }
	}
}

// updateStakingPool reads from the LCD the current staking pool and stores its value inside the database
func (m *Module) updateStakingPool(ctx context.Context, height int64, stakingClient stakingtypes.QueryClient) error {
	// log.Debug().Str("module", "staking").Int64("height", height).
	//	Msg("updating staking pool")

	pool, err := stakingutils.GetStakingPool(height, stakingClient)
	if err != nil {
		return errors.Wrap(err, "GetStakingPool error")
	}

	// TODO: test IT
	if err = m.broker.PublishStakingPool(ctx, m.tbM.MapStakingPool(pool)); err != nil {
		return errors.Wrap(err, "PublishStakingPool error")
	}
	return nil
}

// updateElapsedDelegations updates the redelegations and unbonding delegations that have elapsed
func (m *Module) updateElapsedDelegations(
	ctx context.Context, height int64, timestamp time.Time, enabledModules []string,
) error {
	// log.Debug().Str("module", "staking").Int64("height", height).
	//	Msg("updating elapsed redelegations and unbonding delegations")
	//
	// deletedRedelegations, err := db.DeleteCompletedRedelegations(timestamp)
	// if err != nil {
	//	log.Error().Str("module", "staking").Err(err).Int64("height", height).
	//		Msg("error while deleting completed redelegations")
	//	return
	// }

	// deletedUnbondingDelegations, err := db.DeleteCompletedUnbondingDelegations(timestamp)
	// if err != nil {
	//	log.Error().Str("module", "staking").Err(err).Int64("height", height).
	//		Msg("error while deleting completed unbonding delegations")
	//	return
	// }

	var delegators = map[string]bool{}

	// Add all the delegators from the redelegations
	// for _, redelegation := range deletedRedelegations {
	//	if _, ok := delegators[redelegation.DelegatorAddress]; !ok {
	//		delegators[redelegation.DelegatorAddress] = true
	//	}
	// }
	//
	//// Add all the delegators from unbonding delegations
	// for _, delegation := range deletedUnbondingDelegations {
	//	if _, ok := delegators[delegation.DelegatorAddress]; !ok {
	//		delegators[delegation.DelegatorAddress] = true
	//	}
	// }

	// Update the delegations and balances of all the delegators
	for delegator := range delegators {
		stakingutils.RefreshDelegations(ctx, height, delegator, m.client.StakingQueryClient, m.broker, m.tbM)

		// TODO
		// stakingutils.RefreshBalance(delegator, m.client.BankQueryClient)

		if utils.ContainAny(enabledModules, "history") {
			// TODO:
			// err := historyutils.UpdateAccountBalanceHistory(delegator)
			// if err != nil {
			//	log.Error().Str("module", "staking").Err(err).Int64("height", height).
			//		Str("account", delegator).
			//		Msg("error while updating account balance history")
			//	return
			// }
		}
	}

	return nil
}
