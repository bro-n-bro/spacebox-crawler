package utils

import (
	"context"

	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	tb "github.com/hexy-dev/spacebox-crawler/pkg/mapper/to_broker"
	"github.com/hexy-dev/spacebox-crawler/types"
	"github.com/hexy-dev/spacebox/broker/model"
)

// UpdateDelegationsAndReplaceExisting updates the delegations of the given delegator by querying them at the
// required height, and then publishes them to the broker by replacing all existing ones.
func UpdateDelegationsAndReplaceExisting(
	ctx context.Context,
	height int64,
	delegator string,
	client stakingtypes.QueryClient,
	mapper tb.ToBroker,
	broker interface {
		PublishDelegation(ctx context.Context, d model.Delegation) error
	},
) error {
	// TODO:
	// Remove existing delegations
	// err := db.DeleteDelegatorDelegations(delegator)
	// if err != nil {
	//	return err
	// }

	// Get the delegations
	respPb, err := client.DelegatorDelegations(
		ctx,
		&stakingtypes.QueryDelegatorDelegationsRequest{
			DelegatorAddr: delegator,
		},
	)
	if err != nil {
		return err
	}

	for _, delegation := range respPb.DelegationResponses {
		// TODO: test IT
		if err = broker.PublishDelegation(ctx, model.Delegation{
			OperatorAddress:  delegation.Delegation.ValidatorAddress,
			DelegatorAddress: delegation.Delegation.DelegatorAddress,
			Height:           height,
			Coin:             mapper.MapCoin(types.NewCoinFromCdk(delegation.Balance)),
		}); err != nil {
			return err
		}
	}

	return err
}
