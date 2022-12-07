package core

import (
	"context"

	"bro-n-bro-osmosis/types"
)

func (m *Module) HandleTx(ctx context.Context, tx *types.Tx) error {
	if err := m.broker.PublishTransaction(ctx, m.tbM.MapTransaction(tx)); err != nil {
		m.log.Error().Err(err).Int64("height", tx.Height).Msg("PublishTransaction error")
		return err
	}
	return nil
}
