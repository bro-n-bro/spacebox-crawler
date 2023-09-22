package auth

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/bro-n-bro/spacebox-crawler/modules/utils"
	"github.com/bro-n-bro/spacebox-crawler/types"
)

func (m *Module) HandleMessage(ctx context.Context, _ int, msg sdk.Msg, tx *types.Tx) error {
	addresses := m.parser(m.cdc, msg)

	for _, addr := range addresses {
		if err := utils.GetAndPublishAccount(ctx, addr, tx.Height, m.accCache, m.broker,
			m.client.AuthQueryClient); err != nil {
			return err
		}
	}

	return nil
}
