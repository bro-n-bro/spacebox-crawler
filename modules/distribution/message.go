package distribution

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	"google.golang.org/grpc/codes"

	grpcClient "github.com/hexy-dev/spacebox-crawler/client/grpc"
	"github.com/hexy-dev/spacebox-crawler/types"
	"github.com/hexy-dev/spacebox/broker/model"
)

const (
	errDelegationDoesNotExists = `rpc error: code = %s desc = delegation does not exist`
)

func (m *Module) HandleMessage(ctx context.Context, index int, cosmosMsg sdk.Msg, tx *types.Tx) error {
	if len(tx.Logs) == 0 {
		return nil
	}

	switch msg := cosmosMsg.(type) {
	// TODO: todo to handle block
	case *distrtypes.MsgFundCommunityPool:
		res, err := m.client.DistributionQueryClient.CommunityPool(ctx, &distrtypes.QueryCommunityPoolRequest{})
		if err != nil {
			return err
		}

		// TODO: test it
		return m.broker.PublishCommunityPool(ctx, model.CommunityPool{
			Height: tx.Height,
			Coins:  m.tbM.MapCoins(types.NewCoinsFromCdkDec(res.Pool)),
		})
	case *distrtypes.MsgWithdrawDelegatorReward:
		resp, err := m.client.DistributionQueryClient.DelegationRewards(
			ctx,
			&distrtypes.QueryDelegationRewardsRequest{
				DelegatorAddress: msg.DelegatorAddress,
				ValidatorAddress: msg.ValidatorAddress,
			},
			grpcClient.GetHeightRequestHeader(tx.Height))
		if err != nil {
			// Get the error code
			var code string
			if _, err = fmt.Sscanf(err.Error(), errDelegationDoesNotExists, &code); err != nil {
				return err
			}

			if code == codes.Unknown.String() {
				return nil
			}

			return err
		}

		// TODO: test it
		return m.broker.PublishDelegationRewardMessage(ctx, model.DelegationRewardMessage{
			Coins:            m.tbM.MapCoins(types.NewCoinsFromCdkDec(resp.Rewards)),
			Height:           tx.Height,
			DelegatorAddress: msg.DelegatorAddress,
			ValidatorAddress: msg.ValidatorAddress,
			TxHash:           tx.TxHash,
			MsgIndex:         int64(index),
		})
	}

	return nil
}
