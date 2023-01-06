package distribution

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	"github.com/hexy-dev/spacebox/broker/model"

	grpcClient "github.com/hexy-dev/spacebox-crawler/client/grpc"
	"github.com/hexy-dev/spacebox-crawler/types"
)

func (m *Module) HandleMessage(ctx context.Context, _ int, cosmosMsg sdk.Msg, tx *types.Tx) error {
	if len(tx.Logs) == 0 { // TODO: maybe not needed
		return nil
	}

	switch msg := cosmosMsg.(type) {
	// TODO: todo to handle block
	case *distrtypes.MsgFundCommunityPool:
		res, err := m.client.DistributionQueryClient.CommunityPool(ctx, &distrtypes.QueryCommunityPoolRequest{})
		if err != nil {
			return err
		}

		pool := model.NewCommunityPool(tx.Height, m.tbM.MapCoins(types.NewCoinsFromCdkDec(res.Pool)))
		// TODO: test it
		return m.broker.PublishCommunityPool(ctx, pool)
	case *distrtypes.MsgWithdrawDelegatorReward:
		resp, err := m.client.DistributionQueryClient.DelegationRewards(
			ctx,
			&distrtypes.QueryDelegationRewardsRequest{
				DelegatorAddress: msg.DelegatorAddress,
				ValidatorAddress: msg.ValidatorAddress,
			},
			grpcClient.GetHeightRequestHeader(tx.Height))
		if err != nil {
			return err
		}

		// TODO: question in MIRO
		// if err := m.broker.PublishDelegationReward(ctx, model.NewDelegationReward(
		//	tx.Height, msg.DelegatorAddress, msg.ValidatorAddress,
		//	m.tbM.MapCoins(types.NewCoinsFromCdkDec(resp.Rewards)))); err != nil {
		//	return err
		// }

		// TODO: test it
		return m.broker.PublishDelegationRewardMessage(ctx, model.NewDelegationRewardMessage(
			tx.Height, msg.DelegatorAddress, msg.ValidatorAddress,
			tx.TxHash, m.tbM.MapCoins(types.NewCoinsFromCdkDec(resp.Rewards))))
	}

	return nil
}
