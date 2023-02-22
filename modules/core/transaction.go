package core

import (
	"context"

	"github.com/bro-n-bro/spacebox-crawler/types"
)

func (m *Module) HandleTx(ctx context.Context, tx *types.Tx) error {
	mappedTx, err := m.tbM.MapTransaction(tx)
	if err != nil {
		return err
	}

	if err = m.broker.PublishTransaction(ctx, mappedTx); err != nil {
		m.log.Error().Err(err).Int64("height", tx.Height).Msg("PublishTransaction error")
		return err
	}

	return nil
}
