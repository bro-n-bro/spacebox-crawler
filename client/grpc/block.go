package grpc

import (
	"context"

	"github.com/cosmos/cosmos-sdk/client/grpc/tmservice"
	tmccoretypes "github.com/tendermint/tendermint/rpc/core/types"
	tmctypes "github.com/tendermint/tendermint/types"
)

func (c *Client) Block(ctx context.Context, height int64) (*tmccoretypes.ResultBlock, error) {
	resp, err := c.TmsService.GetBlockByHeight(
		ctx,
		&tmservice.GetBlockByHeightRequest{
			Height: height,
		},
	)
	if err != nil {
		return nil, err
	}

	block, err := tmctypes.BlockFromProto(resp.Block) // nolint:staticcheck
	if err != nil {
		return nil, err
	}

	blockID, err := tmctypes.BlockIDFromProto(resp.BlockId)
	if err != nil {
		return nil, err
	}

	return &tmccoretypes.ResultBlock{
		Block:   block,
		BlockID: *blockID,
	}, nil
}
