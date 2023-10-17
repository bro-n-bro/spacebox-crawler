package rpc

import (
	"context"

	abci "github.com/cometbft/cometbft/abci/types"
)

func (c *Client) GetTxResults(ctx context.Context, height int64) ([]*abci.ResponseDeliverTx, error) {
	result, err := c.RPCClient.BlockResults(ctx, &height)
	if err != nil {
		return nil, err
	}

	return result.TxsResults, nil
}
