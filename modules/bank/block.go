package bank

import (
	"context"

	tmctypes "github.com/tendermint/tendermint/rpc/core/types"

	"bro-n-bro-osmosis/types"
)

func (m *Module) HandleBlock(ctx context.Context, block *types.Block, _ *tmctypes.ResultValidators) error {
	coins, err := m.client.GetTotalSupply(ctx, block.Height)
	if err != nil {
		return err
	}
	_ = coins
	// TODO:
	//err = m.broker.PublishBank(ctx, coins)
	return err
}
