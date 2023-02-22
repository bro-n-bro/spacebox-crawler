package grpc

import (
	"context"

	basetypes "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
)

func (c *Client) GetTotalSupply(ctx context.Context, height int64) (basetypes.Coins, error) {
	respPb, err := c.BankQueryClient.TotalSupply(
		ctx,
		&banktypes.QueryTotalSupplyRequest{},
		GetHeightRequestHeader(height),
	)
	if err != nil {
		return nil, err
	}

	return respPb.Supply, nil
}
