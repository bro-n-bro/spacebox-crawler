package distribution

import (
	"context"

	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	"github.com/pkg/errors"

	grpcClient "github.com/bro-n-bro/spacebox-crawler/client/grpc"
	"github.com/bro-n-bro/spacebox-crawler/types"
	"github.com/bro-n-bro/spacebox/broker/model"
)

func (m *Module) HandleBlock(ctx context.Context, block *types.Block) error {
	// TODO: maybe use consensus client for get correct validators?

	// g, ctx := errgroup.WithContext(ctx)
	// g.Go(func() error {
	//	return m.updateParams(ctx, block.Height)
	// })

	/*
		TODO:
			UpdateValidatorsCommissionAmounts,
			UpdateDelegatorsRewardsAmounts (we need info from some storage for these actions)
	*/

	// if err := g.Wait(); err != nil {
	//	return err
	// }

	res, err := m.client.DistributionQueryClient.CommunityPool(ctx, &distrtypes.QueryCommunityPoolRequest{})
	if err != nil {
		return errors.Wrap(err, "get CommunityPool error")
	}

	// TODO: test it
	if err := m.broker.PublishCommunityPool(ctx, model.CommunityPool{
		Height: block.Height,
		Coins:  m.tbM.MapCoins(types.NewCoinsFromCdkDec(res.Pool)),
	}); err != nil {
		return errors.Wrap(err, "publish CommunityPool error")
	}

	return m.updateParams(ctx, block.Height)
}

func (m *Module) updateParams(ctx context.Context, height int64) error {
	res, err := m.client.DistributionQueryClient.Params(
		ctx,
		&distrtypes.QueryParamsRequest{},
		grpcClient.GetHeightRequestHeader(height),
	)
	if err != nil {
		return err
	}

	if err := m.broker.PublishDistributionParams(ctx, model.DistributionParams{
		Height: height,
		Params: model.DParams{
			CommunityTax:        res.Params.CommunityTax.MustFloat64(),
			BaseProposerReward:  res.Params.BaseProposerReward.MustFloat64(),  //nolint:staticcheck
			BonusProposerReward: res.Params.BonusProposerReward.MustFloat64(), //nolint:staticcheck
			WithdrawAddrEnabled: res.Params.WithdrawAddrEnabled,
		},
	}); err != nil {
		m.log.Error().Int64("height", height).Err(err).Msg("PublishDistributionParams error")
		return err
	}

	return nil
}
