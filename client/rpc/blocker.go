package rpc

import (
	"context"
	"time"

	"github.com/bro-n-bro/spacebox-crawler/types"
)

// GetBlockEvents returns begin block and end block events.
func (c *Client) GetBlockEvents(ctx context.Context, height int64) (begin, end types.BlockerEvents, err error) {
	ctx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	result, err := c.RPCClient.BlockResults(ctx, &height)
	if err != nil {
		return nil, nil, err
	}

	begin = types.NewBlockerEventsAttributes(result.BeginBlockEvents)
	end = types.NewBlockerEventsAttributes(result.EndBlockEvents)

	return
}
