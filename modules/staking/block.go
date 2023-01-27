package staking

import (
	"context"
	"fmt"

	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"golang.org/x/sync/errgroup"

	grpcClient "github.com/bro-n-bro/spacebox-crawler/client/grpc"
	"github.com/bro-n-bro/spacebox-crawler/types"
	"github.com/bro-n-bro/spacebox/broker/model"
)

func (m *Module) HandleBlock(ctx context.Context, block *types.Block) error {
	g, ctx2 := errgroup.WithContext(ctx)

	g.Go(func() error {
		// Update the params
		return m.updateParams(ctx2, block.Height)
	})

	g.Go(func() error {
		// Update the staking pool
		return m.updateStakingPool(ctx2, block.Height, m.client.StakingQueryClient)
	})

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
			MaxValidators:     uint64(res.Params.MaxValidators),
			MaxEntries:        uint64(res.Params.MaxEntries),
			HistoricalEntries: uint64(res.Params.HistoricalEntries),
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
