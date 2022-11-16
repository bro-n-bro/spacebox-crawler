package utils

import (
	"context"

	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	grpcClient "bro-n-bro-osmosis/client/grpc"
	"bro-n-bro-osmosis/types"
)

func GetStakingPool(height int64, stakingClient stakingtypes.QueryClient) (*types.Pool, error) {
	res, err := stakingClient.Pool(
		context.Background(),
		&stakingtypes.QueryPoolRequest{},
		grpcClient.GetHeightRequestHeader(height),
	)
	if err != nil {
		return nil, err
	}

	return types.NewPool(res.Pool.BondedTokens, res.Pool.NotBondedTokens, height), nil
}
