package bank

import (
	"context"

	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	grpcClient "github.com/hexy-dev/spacebox-crawler/client/grpc"
	tmctypes "github.com/tendermint/tendermint/rpc/core/types"

	"github.com/hexy-dev/spacebox-crawler/types"
)

func (m *Module) HandleBlock(ctx context.Context, block *types.Block, _ *tmctypes.ResultValidators) error {
	resp, err := m.client.BankQueryClient.TotalSupply(
		ctx,
		&banktypes.QueryTotalSupplyRequest{},
		grpcClient.GetHeightRequestHeader(block.Height))
	if err != nil {
		return err
	}

	// TODO: tests
	err = m.broker.PublishSupply(ctx, m.tbM.MapSupply(types.NewTotalSupply(block.Height, types.NewCoinsFromCdk(resp.Supply))))
	if err != nil {
		return err
	}
	return nil
}
