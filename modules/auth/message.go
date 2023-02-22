package auth

import (
	"context"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/bro-n-bro/spacebox-crawler/types"
	"github.com/bro-n-bro/spacebox/broker/model"
)

func (m *Module) HandleMessage(ctx context.Context, _ int, msg sdk.Msg, tx *types.Tx) error {
	addresses := m.parser(m.cdc, msg)
	prefix := sdk.GetConfig().GetBech32AccountAddrPrefix()

	for _, addr := range addresses {
		if strings.HasPrefix(addr, prefix) {
			// TODO: test it
			if err := m.broker.PublishAccount(ctx, model.Account{
				Address: addr,
				Height:  tx.Height,
			}); err != nil {
				return err
			}
		}
	}

	return nil
}
