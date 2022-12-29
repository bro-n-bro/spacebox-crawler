package auth

import (
	"context"

	"github.com/hexy-dev/spacebox-crawler/modules/utils"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/hexy-dev/spacebox-crawler/types"
)

func (m *Module) HandleMessage(ctx context.Context, _ int, msg sdk.Msg, tx *types.Tx) error {
	addresses, err := m.parser(m.cdc, msg)
	if err != nil {
		m.log.Error().Err(err).Msg("HandleMessage getAddresses error")
		return nil
	}

	// TODO: test it
	err = m.broker.PublishAccounts(ctx, m.tbM.MapAccounts(utils.GetAccounts(addresses, tx.Height)))
	if err != nil {
		return err
	}
	return nil
}
