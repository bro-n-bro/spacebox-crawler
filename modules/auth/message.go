package auth

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"bro-n-bro-osmosis/types"
)

func (m *Module) HandleMessage(ctx context.Context, index int, msg sdk.Msg, tx *types.Tx) error {
	addresses, err := m.parser(m.cdc, msg)
	if err != nil {
		m.log.Error().Err(err).Msg("HandleMessage getAddresses error")
		return nil
	}

	// TODO:
	_ = addresses
	return nil
}
