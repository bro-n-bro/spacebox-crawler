package rpc

import (
	"context"
	"time"

	abci "github.com/cometbft/cometbft/abci/types"
)

func (c *Client) GetTxResults(ctx context.Context, height int64) ([]*abci.ResponseDeliverTx, error) {
	ctx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	result, err := c.RPCClient.BlockResults(ctx, &height)
	if err != nil {
		return nil, err
	}

	return result.TxsResults, nil
}
