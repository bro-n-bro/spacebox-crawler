package utils

import (
	"context"

	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	grpcClient "github.com/hexy-dev/spacebox-crawler/client/grpc"
	"github.com/hexy-dev/spacebox-crawler/types"
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
