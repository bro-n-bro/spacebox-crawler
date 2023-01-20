package auth

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/bro-n-bro/spacebox-crawler/types"
	"github.com/bro-n-bro/spacebox/broker/model"
)

func (m *Module) HandleMessage(ctx context.Context, _ int, msg sdk.Msg, tx *types.Tx) error {
	addresses, err := m.parser(m.cdc, msg)
	if err != nil {
		m.log.Error().Err(err).Msg("HandleMessage getAddresses error")
		return nil
	}

	for _, addr := range addresses {
		// TODO: test it
		if err = m.broker.PublishAccount(ctx, model.Account{
			Address: addr,
			Height:  tx.Height,
		}); err != nil {
			return err
		}
	}

	return nil
}
