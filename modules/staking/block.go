package staking

import (
	"context"
	"fmt"

	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"golang.org/x/sync/errgroup"

	grpcClient "github.com/hexy-dev/spacebox-crawler/client/grpc"
	stakingutils "github.com/hexy-dev/spacebox-crawler/modules/staking/utils"
	"github.com/hexy-dev/spacebox-crawler/types"
	"github.com/hexy-dev/spacebox/broker/model"
)

func (m *Module) HandleBlock(ctx context.Context, block *types.Block) error {
	// Update the validators
	validators, err := stakingutils.UpdateValidators(ctx, block.Height, m.client.StakingQueryClient, m.cdc, m.broker)
	if err != nil {
		return err
	}

	g, ctx2 := errgroup.WithContext(ctx)

	g.Go(func() error {
		// Update the params
		return m.updateParams(ctx2, block.Height)
	})

	g.Go(func() error {
		// Update the validators statuses
		return m.updateValidatorsStatus(ctx2, block.Height, validators)
	})

	g.Go(func() error {
		// Update the staking pool
		return m.updateStakingPool(ctx2, block.Height, m.client.StakingQueryClient)
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

	return g.Wait()
}

// updateParams gets the updated params and stores them inside the database
func (m *Module) updateParams(ctx context.Context, height int64) error {
	res, err := m.client.StakingQueryClient.Params(
		ctx,
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

	// TODO: test it
	// TODO: maybe check diff from mongo in my side?
	if err = m.broker.PublishStakingParams(ctx, model.StakingParams{
		Height: height,
		Params: model.SParams{
			UnbondingTime:     res.Params.UnbondingTime,
			MaxValidators:     res.Params.MaxValidators,
			MaxEntries:        res.Params.MaxEntries,
			HistoricalEntries: res.Params.HistoricalEntries,
			BondDenom:         res.Params.BondDenom,
			MinCommissionRate: commissionRate,
		},
	}); err != nil {
		return err
	}

	// TODO:
	// err = db.SaveStakingParams(types.NewStakingParams(res.Params, height))
	// if err != nil {
	// log.Error().Str("module", "staking").Err(err).
	//	Int64("height", height).
	//	Msg("error while saving params")
	// return
	// }

	return nil
}

// updateValidatorsStatus updates all validators' statuses
func (m *Module) updateValidatorsStatus(ctx context.Context, height int64, vals []stakingtypes.Validator) error {
	for _, validator := range vals {
		consAddr, err := stakingutils.GetValidatorConsAddr(m.cdc, validator)
		if err != nil {
			return fmt.Errorf("error while getting validator consensus address: %w", err)
		}

		consPubKey, err := stakingutils.GetValidatorConsPubKey(m.cdc, validator)
		if err != nil {
			return fmt.Errorf("error while getting validator consensus public key: %w", err)
		}

		// TODO: save to mongo?
		// TODO: test it
		if err = m.broker.PublishValidator(ctx, model.Validator{
			ConsensusAddress: consAddr.String(),
			ConsensusPubkey:  consPubKey.String(),
		}); err != nil {
			return err
		}

		// TODO: test it
		if err = m.broker.PublishValidatorStatus(ctx, model.ValidatorStatus{
			Height:           height,
			ValidatorAddress: consAddr.String(),
			Status:           int64(validator.GetStatus()),
			Jailed:           validator.IsJailed(),
		}); err != nil {
			return err
		}
	}

	return nil
}

// updateStakingPool reads from the LCD the current staking pool and stores its value inside the database
func (m *Module) updateStakingPool(ctx context.Context, height int64, stakingClient stakingtypes.QueryClient) error {
	respPb, err := stakingClient.Pool(
		ctx,
		&stakingtypes.QueryPoolRequest{},
		grpcClient.GetHeightRequestHeader(height),
	)
	if err != nil {
		return err
	}

	// TODO: to mapper?
	var bondedTokens, notBondedTokens int64

	if !respPb.Pool.BondedTokens.IsNil() {
		bondedTokens = respPb.Pool.BondedTokens.Int64()
	}

	if !respPb.Pool.NotBondedTokens.IsNil() {
		notBondedTokens = respPb.Pool.NotBondedTokens.Int64()
	}

	// TODO: test IT
	if err = m.broker.PublishStakingPool(ctx, model.StakingPool{
		Height:          height,
		NotBondedTokens: notBondedTokens,
		BondedTokens:    bondedTokens,
	}); err != nil {
		return fmt.Errorf("PublishStakingPool error: %w", err)
	}

	return nil
}
