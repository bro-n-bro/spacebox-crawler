package utils

import (
	"context"
	"fmt"

	grpcClient "github.com/hexy-dev/spacebox-crawler/client/grpc"
	"github.com/hexy-dev/spacebox-crawler/types"

	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
)

// UpdateDelegatorsRewardsAmounts updates the delegators commission amounts
func UpdateDelegatorsRewardsAmounts(height int64, client distrtypes.QueryClient) error {
	// Get the delegators
	// delegators, err := db.GetDelegators()
	// if err != nil {
	// log.Error().Str("module", "distribution").Err(err).Int64("height", height).
	//	Msg("error while getting delegators")
	// }

	var delegators []string
	if len(delegators) == 0 {
		return nil
	}

	// Get the rewards
	for _, delegator := range delegators {
		if err := updateDelegatorCommission(height, delegator, client); err != nil {
			return err
		}
	}
	return nil
}

func updateDelegatorCommission(height int64, delegator string, distrClient distrtypes.QueryClient) error {
	rewards, err := GetDelegatorRewards(height, delegator, distrClient)
	if err != nil {
		return err
	}

	// TODO:
	_ = rewards
	// err = db.SaveDelegatorsRewardsAmounts(rewards)
	// if err != nil {
	//	log.Error().Str("module", "distribution").Err(err).
	//		Int64("height", height).Str("delegator", delegator).
	//		Msg("error while saving delegator rewards")
	// }
	return nil
}

// GetDelegatorRewards returns the current rewards for the given delegator
func GetDelegatorRewards(height int64, delegator string,
	distrClient distrtypes.QueryClient) ([]types.DelegatorRewardAmount, error) {

	header := grpcClient.GetHeightRequestHeader(height)
	rewardsRes, err := distrClient.DelegationTotalRewards(
		context.Background(),
		&distrtypes.QueryDelegationTotalRewardsRequest{DelegatorAddress: delegator},
		header,
	)
	if err != nil {
		return nil, fmt.Errorf("error while getting delegator reward: %s", err)
	}

	withdrawAddressRes, err := distrClient.DelegatorWithdrawAddress(
		context.Background(),
		&distrtypes.QueryDelegatorWithdrawAddressRequest{DelegatorAddress: delegator},
		header,
	)
	if err != nil {
		return nil, fmt.Errorf("error while getting delegator rewards: %s", err)
	}

	var rewards = make([]types.DelegatorRewardAmount, len(rewardsRes.Rewards))
	for index, reward := range rewardsRes.Rewards {
		rewards[index] = types.NewDelegatorRewardAmount(
			delegator,
			reward.ValidatorAddress,
			withdrawAddressRes.WithdrawAddress,
			reward.Reward,
			height,
		)
	}
	return rewards, nil
}
