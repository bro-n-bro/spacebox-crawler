package rpc

import (
	"context"
	"time"

	tmctypes "github.com/tendermint/tendermint/rpc/core/types"
)

func (c *Client) SubscribeNewBlocks(subscriber string) (<-chan tmctypes.ResultEvent, context.CancelFunc, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	eventCh, err := c.RpcClient.Subscribe(ctx, subscriber, "tm.event = 'NewBlock'")
	return eventCh, cancel, err
}
