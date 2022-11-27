package utils

import (
	"context"

	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
)

// UpdateCommunityPool fetch total amount of coins in the system from RPC and store it into database
func UpdateCommunityPool(ctx context.Context, height int64, client distrtypes.QueryClient) error {
	//log.Debug().Str("module", "distribution").Int64("height", height).Msg("getting community pool")

	res, err := client.CommunityPool(ctx, &distrtypes.QueryCommunityPoolRequest{})
	if err != nil {
		return err
	}

	// TODO:
	_ = res

	// Store the signing infos into the database
	//return db.SaveCommunityPool(res.Pool, height)
	return nil
}
