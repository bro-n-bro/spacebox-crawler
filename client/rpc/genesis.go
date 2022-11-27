package rpc

import (
	"context"

	tmtypes "github.com/tendermint/tendermint/types"
)

func (c *Client) Genesis(ctx context.Context) (*tmtypes.GenesisDoc, error) {
	// FIXME
	g, err := c.RpcClient.Genesis(ctx)
	if err != nil {
		return nil, err
	}
	return g.Genesis, nil
}
