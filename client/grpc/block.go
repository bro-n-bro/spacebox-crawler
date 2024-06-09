package grpc

import (
	"context"

	"github.com/cosmos/cosmos-sdk/client/grpc/tmservice"
	coretypes "github.com/tendermint/tendermint/rpc/core/types"
	tmtypes "github.com/tendermint/tendermint/types"
)

func (c *Client) Block(ctx context.Context, height int64) (*coretypes.ResultBlock, error) {
	resp, err := c.TmsService.GetBlockByHeight(
		ctx,
		&tmservice.GetBlockByHeightRequest{
			Height: height,
		},
	)
	if err != nil {
		return nil, err
	}

	block, err := tmtypes.BlockFromProto(resp.Block) // nolint:staticcheck
	if err != nil {
		return nil, err
	}

	blockID, err := tmtypes.BlockIDFromProto(resp.BlockId)
	if err != nil {
		return nil, err
	}

	return &coretypes.ResultBlock{
		Block:   block,
		BlockID: *blockID,
	}, nil
}
