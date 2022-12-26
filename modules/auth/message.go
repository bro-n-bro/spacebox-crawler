package auth

import (
	"bro-n-bro-osmosis/modules/utils"
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"bro-n-bro-osmosis/types"
)

func (m *Module) HandleMessage(ctx context.Context, _ int, msg sdk.Msg, tx *types.Tx) error {
	addresses, err := m.parser(m.cdc, msg)
	if err != nil {
		m.log.Error().Err(err).Msg("HandleMessage getAddresses error")
		return nil
	}

	// TODO:
	err = m.broker.PublishAccounts(ctx, m.tbM.MapAccounts(utils.GetAccounts(addresses, tx.Height)))
	if err != nil {
		return err
	}
	return nil
}
