package distribution

import (
	"context"

	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	tmctypes "github.com/tendermint/tendermint/rpc/core/types"
	"golang.org/x/sync/errgroup"

	grpcClient "github.com/hexy-dev/spacebox-crawler/client/grpc"
	"github.com/hexy-dev/spacebox-crawler/modules/distribution/utils"
	"github.com/hexy-dev/spacebox-crawler/types"
)

func (m *Module) HandleBlock(ctx context.Context, block *types.Block, vals *tmctypes.ResultValidators) error {
	// TODO: maybe use consensus client for get correct validators?

	g, ctx := errgroup.WithContext(ctx)
	g.Go(func() error {
		return m.updateParams(ctx, block.Height)
	})

	// Update the validator commissions
	g.Go(func() error {
		return utils.UpdateValidatorsCommissionAmounts(block.Height, m.client.DistributionQueryClient)
	})

	// Update the delegators commissions amounts
	g.Go(func() error {
		return utils.UpdateDelegatorsRewardsAmounts(block.Height, m.client.DistributionQueryClient)
	})

	if err := g.Wait(); err != nil {
		return err
	}

	// TODO: client.community pull
	return nil
}

func (m *Module) updateParams(ctx context.Context, height int64) error {
	// log.Debug().Str("module", "distribution").Int64("height", height).
	//	Msg("updating params")

	res, err := m.client.DistributionQueryClient.Params(
		context.Background(),
		&distrtypes.QueryParamsRequest{},
		grpcClient.GetHeightRequestHeader(height),
	)
	if err != nil {
		return err
	}

	// TODO: maybe check diff from mongo in my side?
	params := types.NewDistributionParams(res.Params, height)
	// TODO: test it
	if err := m.broker.PublishDistributionParams(ctx, m.tbM.MapDistributionParams(params)); err != nil {
		m.log.Error().Int64("height", height).Err(err).Msg("PublishDistributionParams error")
		return err
	}
	return nil
}
