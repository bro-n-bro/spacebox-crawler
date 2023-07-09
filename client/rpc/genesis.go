package rpc

import (
	"context"

	cometbfttypes "github.com/cometbft/cometbft/types"
)

func (c *Client) Genesis(ctx context.Context) (*cometbfttypes.GenesisDoc, error) {
	g, err := c.RPCClient.Genesis(ctx)
	if err != nil {
		return nil, err
	}

	return g.Genesis, nil
}
