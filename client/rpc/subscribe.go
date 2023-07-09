package rpc

import (
	"context"

	cometbftcoretypes "github.com/cometbft/cometbft/rpc/core/types"
)

func (c *Client) SubscribeNewBlocks(ctx context.Context) (<-chan cometbftcoretypes.ResultEvent, error) {
	return c.RPCClient.Subscribe(ctx, "", "tm.event = 'NewBlock'")
}
