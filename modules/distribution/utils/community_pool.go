package utils

import (
	"context"

	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"

	"bro-n-bro-osmosis/internal/rep"
	tb "bro-n-bro-osmosis/pkg/mapper/to_broker"
	"bro-n-bro-osmosis/types"
)

// UpdateCommunityPool fetch total amount of coins in the system from RPC and store it into database
func UpdateCommunityPool(ctx context.Context, height int64, client distrtypes.QueryClient, broker rep.Broker,
	mapper tb.ToBroker) error {

	//log.Debug().Str("module", "distribution").Int64("height", height).Msg("getting community pool")

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
