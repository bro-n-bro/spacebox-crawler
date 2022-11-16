package utils

import (
	"context"
	"fmt"
	"sync"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	grpcClient "bro-n-bro-osmosis/client/grpc"
	"bro-n-bro-osmosis/types"
)

// ConvertUnbondingResponse converts the given UnbondingDelegation response into a slice of BDJuno UnbondingDelegation
func ConvertUnbondingResponse(
	height int64, bondDenom string, response stakingtypes.UnbondingDelegation,
) []types.UnbondingDelegation {
	var delegations []types.UnbondingDelegation
	for _, entry := range response.Entries {
		delegations = append(delegations, types.NewUnbondingDelegation(
			response.DelegatorAddress,
			response.ValidatorAddress,
			sdk.NewCoin(bondDenom, entry.Balance),
			entry.CompletionTime,
			height,
		))
	}
	return delegations
}

// --------------------------------------------------------------------------------------------------------------------

// UpdateValidatorsUnbondingDelegations updates the unbonding delegations for all the validators provided
func UpdateValidatorsUnbondingDelegations(
	height int64, bondDenom string, validators []stakingtypes.Validator,
	client stakingtypes.QueryClient,
) {
	var wg sync.WaitGroup
	for _, val := range validators {
		wg.Add(1)
		go getUnbondingDelegations(val.OperatorAddress, bondDenom, height, client, &wg)
	}
	wg.Wait()
}

// getUnbondingDelegations gets all the unbonding delegations referring to the validator having the
// given address at the given block height (having the given timestamp).
// All the unbonding delegations will be sent to the out channel, and wg.Done() will be called at the end.
func getUnbondingDelegations(
	validatorAddress string, bondDenom string, height int64,
	stakingClient stakingtypes.QueryClient, wg *sync.WaitGroup,
) {
	defer wg.Done()

	header := grpcClient.GetHeightRequestHeader(height)

	var nextKey []byte
	var stop = false
	for !stop {
		res, err := stakingClient.ValidatorUnbondingDelegations(
			context.Background(),
			&stakingtypes.QueryValidatorUnbondingDelegationsRequest{
				ValidatorAddr: validatorAddress,
				Pagination: &query.PageRequest{
					Key:   nextKey,
					Limit: 100, // Query 100 unbonding delegations at a time
				},
			},
			header,
		)
		if err != nil {
			return
		}

		var delegations []types.UnbondingDelegation
		for _, delegation := range res.UnbondingResponses {
			delegations = append(delegations, ConvertUnbondingResponse(height, bondDenom, delegation)...)
		}

		// TODO:
		//err = db.SaveUnbondingDelegations(delegations)
		//if err != nil {
		//	return
		//}

		nextKey = res.Pagination.NextKey
		stop = len(res.Pagination.NextKey) == 0
	}
}

// GetDelegatorUnbondingDelegations returns the current unbonding delegations for the user having the given address
func GetDelegatorUnbondingDelegations(
	height int64, address string, bondDenom string, stakingClient stakingtypes.QueryClient,
) ([]types.UnbondingDelegation, error) {
	var delegations []types.UnbondingDelegation

	header := grpcClient.GetHeightRequestHeader(height)

	var nextKey []byte
	var stop = false
	for !stop {
		res, err := stakingClient.DelegatorUnbondingDelegations(
			context.Background(),
			&stakingtypes.QueryDelegatorUnbondingDelegationsRequest{
				DelegatorAddr: address,
				Pagination: &query.PageRequest{
					Key:   nextKey,
					Limit: 100, // Query 100 unbonding delegations at a time
				},
			},
			header,
		)
		if err != nil {
			return nil, fmt.Errorf("error while getting validator delegations: %s", err)
		}

		for _, delegation := range res.UnbondingResponses {
			delegations = append(delegations, ConvertUnbondingResponse(height, bondDenom, delegation)...)
		}

		nextKey = res.Pagination.NextKey
		stop = len(res.Pagination.NextKey) == 0
	}

	return delegations, nil
}

// --------------------------------------------------------------------------------------------------------------------

// DeleteUnbondingDelegation returns a function that when called deletes the given delegation from the database
func DeleteUnbondingDelegation(delegation types.UnbondingDelegation) func() {
	return func() {
		// TODO:
		//err := db.DeleteUnbondingDelegation(delegation)
		//if err != nil {
		//
		//}
	}
}
