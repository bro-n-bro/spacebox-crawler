package bank

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"bro-n-bro-osmosis/types"
)

func (m *Module) HandleMessage(_ context.Context, _ int, cosmosMsg sdk.Msg, _ *types.Tx) error {
	addresses, err := m.parser(m.cdc, cosmosMsg)
	if err != nil {
		m.log.Error().Err(err).Msg("HandleMessage getAddresses error:")
		return nil
	}
	//err = m.broker.PublishBank(ctx)

	// todo: publish?
	_ = addresses

	return nil
}
