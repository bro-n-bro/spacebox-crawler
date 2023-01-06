package core

import (
	"context"

	"github.com/hexy-dev/spacebox/broker/model"

	"github.com/hexy-dev/spacebox-crawler/types"
)

func (m *Module) HandleBlock(ctx context.Context, block *types.Block) error {
	if err := m.broker.PublishBlock(ctx, model.NewBlock(
		block.Height,
		block.Hash,
		block.ProposerAddress,
		block.TxNum,
		block.TotalGas,
		block.Timestamp,
	)); err != nil {
		m.log.Error().
			Err(err).
			Int64("block_height", block.Height).
			Msg("PublishBlock error")
		return err
	}

	return nil

}
