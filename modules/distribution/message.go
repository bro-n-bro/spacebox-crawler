package distribution

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"

	grpcClient "bro-n-bro-osmosis/client/grpc"
	"bro-n-bro-osmosis/modules/distribution/utils"
	"bro-n-bro-osmosis/types"
)

func (m *Module) HandleMessage(ctx context.Context, _ int, cosmosMsg sdk.Msg, tx *types.Tx) error {
	if len(tx.Logs) == 0 { // TODO: maybe not needed
		return nil
	}

	switch msg := cosmosMsg.(type) {
	// TODO: todo to handle block
	case *distrtypes.MsgFundCommunityPool:
		return utils.UpdateCommunityPool(ctx, tx.Height, m.client.DistributionQueryClient, m.broker, m.tbM)
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

		drm := types.NewDelegationRewardMessage(
			msg.DelegatorAddress,
			msg.ValidatorAddress,
			tx.TxHash,
			tx.Height,
			types.NewCoinsFromCdkDec(resp.Rewards),
		)

		// TODO: test it
		return m.broker.PublishDelegationRewardMessage(ctx, m.tbM.MapDelegationRewardMessage(drm))
	}

	return nil
}
