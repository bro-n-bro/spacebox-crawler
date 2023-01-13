package rpc

import (
	"context"

	tmctypes "github.com/tendermint/tendermint/rpc/core/types"
)

func (c *Client) SubscribeNewBlocks(ctx context.Context) (<-chan tmctypes.ResultEvent, error) {
	return c.RPCClient.Subscribe(ctx, "", "tm.event = 'NewBlock'")
}
