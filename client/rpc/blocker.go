package rpc

import (
	"context"

	abci "github.com/tendermint/tendermint/abci/types"
)

func (c *Client) GetBlockEvents(ctx context.Context, height int64) (begin []abci.Event, end []abci.Event, err error) {
	result, err := c.RPCClient.BlockResults(ctx, &height)
	if err != nil {
		return nil, nil, err
	}
	return result.BeginBlockEvents, result.EndBlockEvents, nil
}
