package bank

import (
	grpcClient "bro-n-bro-osmosis/client/grpc"
	"context"

	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"

	tmctypes "github.com/tendermint/tendermint/rpc/core/types"

	"bro-n-bro-osmosis/types"
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
