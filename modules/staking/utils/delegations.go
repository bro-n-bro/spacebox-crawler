package utils

import (
	"context"
	"sync"

	"github.com/cosmos/cosmos-sdk/types/query"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	grpcClient "bro-n-bro-osmosis/client/grpc"
	"bro-n-bro-osmosis/internal/rep"
	tb "bro-n-bro-osmosis/pkg/mapper/to_broker"
	"bro-n-bro-osmosis/types"
)

// ConvertDelegationResponse converts the given response to a BDJuno Delegation instance
func ConvertDelegationResponse(height int64, response stakingtypes.DelegationResponse) types.Delegation {
	return types.NewDelegation(
		response.Delegation.DelegatorAddress,
		response.Delegation.ValidatorAddress,
		response.Balance,
		height,
	)
}

// --------------------------------------------------------------------------------------------------------------------

// UpdateValidatorsDelegations updates the delegations for all the given validators at the provided height
func UpdateValidatorsDelegations(
	height int64, validators []stakingtypes.Validator, client stakingtypes.QueryClient,
) {
	var wg sync.WaitGroup
	for _, val := range validators {
		wg.Add(1)
		go getDelegationsFromGrpc(val.OperatorAddress, height, client, &wg)
	}
	wg.Wait()
}

// getDelegationsFromGrpc gets the list of all the delegations that the validator having the given address has
// at the given block height (having the given timestamp).
// All the delegations will be sent to the out channel, and wg.Done() will be called at the end.
func getDelegationsFromGrpc(
	validatorAddress string, height int64, stakingClient stakingtypes.QueryClient, wg *sync.WaitGroup,
) {
	defer wg.Done()

	header := grpcClient.GetHeightRequestHeader(height)

	var nextKey []byte
	var stop = false
	for !stop {
		res, err := stakingClient.ValidatorDelegations(
			context.Background(),
			&stakingtypes.QueryValidatorDelegationsRequest{
				ValidatorAddr: validatorAddress,
				Pagination: &query.PageRequest{
					Key:   nextKey,
					Limit: 100, // Query 100 delegations at time
				},
			},
			header,
		)
		if err != nil {

			return
		}

		var delegations = make([]types.Delegation, len(res.DelegationResponses))
		for index, delegation := range res.DelegationResponses {
			delegations[index] = ConvertDelegationResponse(height, delegation)
		}

		// TODO:
		//err = db.SaveDelegations(delegations)
		//if err != nil {
		//	return
		//}

		nextKey = res.Pagination.NextKey
		stop = len(res.Pagination.NextKey) == 0
	}
}

// --------------------------------------------------------------------------------------------------------------------

// GetDelegatorDelegations returns the current delegations for the given delegator
func GetDelegatorDelegations(height int64, delegator string, client stakingtypes.QueryClient) ([]types.Delegation, error) {
	// Get the delegations
	res, err := client.DelegatorDelegations(
		context.Background(),
		&stakingtypes.QueryDelegatorDelegationsRequest{
			DelegatorAddr: delegator,
		},
	)
	if err != nil {
		return nil, err
	}

	var delegations = make([]types.Delegation, len(res.DelegationResponses))
	for index, delegation := range res.DelegationResponses {
		delegations[index] = ConvertDelegationResponse(height, delegation)
	}

	return delegations, nil
}

// --------------------------------------------------------------------------------------------------------------------

// UpdateDelegationsAndReplaceExisting updates the delegations of the given delegator by querying them at the
// required height, and then stores them inside the database by replacing all existing ones.
func UpdateDelegationsAndReplaceExisting(ctx context.Context, height int64, delegator string,
	client stakingtypes.QueryClient, broker rep.Broker, mapper tb.ToBroker) error {
	// Remove existing delegations
	//err := db.DeleteDelegatorDelegations(delegator)
	//if err != nil {
	//	return err
	//}

	delegations, err := GetDelegatorDelegations(height, delegator, client)
	if err != nil {
		return err
	}

	for _, delegation := range delegations {
		// TODO: test IT
		if err = broker.PublishDelegation(ctx, mapper.MapDelegation(delegation)); err != nil {
			return err
		}
	}

	return err
}

// RefreshDelegations returns a function that when called updates the delegations of the provided delegator.
// In order to properly update the data, it removes all the existing delegations and stores new ones querying the gRPC
func RefreshDelegations(ctx context.Context, height int64, delegator string, client stakingtypes.QueryClient,
	broker rep.Broker, mapper tb.ToBroker) func() {
	return func() {
		err := UpdateDelegationsAndReplaceExisting(ctx, height, delegator, client, broker, mapper)
		if err != nil {
			//log.Error().Str("module", "staking").Err(err).
			//	Str("operation", "refresh delegations").Msg("error while refreshing delegations")
		}
	}
}
