package staking

import (
	"context"
	"fmt"

	"github.com/hexy-dev/spacebox/broker/model"

	"cosmossdk.io/errors"

	"golang.org/x/sync/errgroup"

	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	grpcClient "github.com/hexy-dev/spacebox-crawler/client/grpc"
	stakingutils "github.com/hexy-dev/spacebox-crawler/modules/staking/utils"
	"github.com/hexy-dev/spacebox-crawler/types"
)

func (m *Module) HandleBlock(ctx context.Context, block *types.Block) error {
	// Update the validators
	validators, err := stakingutils.UpdateValidators(ctx, block.Height, m.client.StakingQueryClient, m.cdc, m.broker)
	if err != nil {
		return err
	}

	g, _ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		// Update the params
		return m.updateParams(_ctx, block.Height)
	})

	g.Go(func() error {
		// Update the validators statuses
		return m.updateValidatorsStatus(_ctx, block.Height, validators)
	})

	g.Go(func() error {
		// Update the staking pool
		return m.updateStakingPool(_ctx, block.Height, m.client.StakingQueryClient)
	})

	// FIXME: does it needed?
	// Update the voting powers
	// go updateValidatorVotingPower(block.Height, vals)

	// FIXME: does it needed?
	// Updated the double sign evidences
	// go updateDoubleSignEvidence(block.Height, block.Evidence.Evidence)

	// TODO: does it needed?
	// g.Go(func() error {
	// Update redelegations and unbonding delegations
	// return m.updateElapsedDelegations(_ctx, block.Height, block.Timestamp, m.enabledModules)
	// })

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

	// TODO: to mapper?
	var commissionRate float64
	if !res.Params.MinCommissionRate.IsNil() {
		commissionRate = res.Params.MinCommissionRate.MustFloat64()
	}

	modelParams := model.NewStakingParams(height, res.Params.MaxValidators, res.Params.MaxEntries,
		res.Params.HistoricalEntries, res.Params.BondDenom, commissionRate, res.Params.UnbondingTime)

	// TODO: test it
	// TODO: maybe check diff from mongo in my side?
	err = m.broker.PublishStakingParams(ctx, modelParams)
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
func (m *Module) updateValidatorsStatus(ctx context.Context, height int64, vals []stakingtypes.Validator) error {
	var (
		val    model.Validator
		status model.ValidatorStatus
	)

	for _, validator := range vals {
		consAddr, err := stakingutils.GetValidatorConsAddr(m.cdc, validator)
		if err != nil {
			return fmt.Errorf("error while getting validator consensus address: %s", err)
		}

		consPubKey, err := stakingutils.GetValidatorConsPubKey(m.cdc, validator)
		if err != nil {
			return fmt.Errorf("error while getting validator consensus public key: %s", err)
		}

		// TODO: save to mongo?
		// TODO: test it
		val = model.NewValidator(consAddr.String(), consPubKey.String())
		err = m.broker.PublishValidators(ctx, []model.Validator{val})
		if err != nil {
			return err
		}

		status = model.NewValidatorStatus(height, int64(validator.GetStatus()), consAddr.String(), validator.IsJailed())
		// TODO: test it
		if err = m.broker.PublishValidatorsStatuses(ctx, []model.ValidatorStatus{status}); err != nil {
			return err
		}
	}

	return nil
}

// updateStakingPool reads from the LCD the current staking pool and stores its value inside the database
func (m *Module) updateStakingPool(ctx context.Context, height int64, stakingClient stakingtypes.QueryClient) error {
	pbResp, err := stakingClient.Pool(
		context.Background(),
		&stakingtypes.QueryPoolRequest{},
		grpcClient.GetHeightRequestHeader(height),
	)
	if err != nil {
		return err
	}

	// TODO: to mapper?
	var bondedTokens, notBondedTokens int64

	if !pbResp.Pool.BondedTokens.IsNil() {
		bondedTokens = pbResp.Pool.BondedTokens.Int64()
	}

	if !pbResp.Pool.NotBondedTokens.IsNil() {
		notBondedTokens = pbResp.Pool.NotBondedTokens.Int64()
	}

	// TODO: test IT
	if err = m.broker.PublishStakingPool(ctx, model.NewStakingPool(height, notBondedTokens, bondedTokens)); err != nil {
		return errors.Wrap(err, "PublishStakingPool error")
	}
	return nil
}
