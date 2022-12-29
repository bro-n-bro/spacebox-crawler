package utils

import (
	"context"

	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"

	"github.com/hexy-dev/spacebox-crawler/internal/rep"
	tb "github.com/hexy-dev/spacebox-crawler/pkg/mapper/to_broker"
	"github.com/hexy-dev/spacebox-crawler/types"
)

// UpdateCommunityPool fetch total amount of coins in the system from RPC and store it into database
func UpdateCommunityPool(ctx context.Context, height int64, client distrtypes.QueryClient, broker rep.Broker,
	mapper tb.ToBroker) error {

	res, err := client.CommunityPool(ctx, &distrtypes.QueryCommunityPoolRequest{})
	if err != nil {
		return err
	}

	// TODO: test it
	err = broker.PublishCommunityPool(ctx, mapper.MapCommunityPool(types.NewCommunityPool(height, res.Pool)))
	if err != nil {
		return err
	}

	return nil
}
