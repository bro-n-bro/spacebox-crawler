package core

import (
	"context"

	"github.com/pkg/errors"

	"github.com/bro-n-bro/spacebox-crawler/types"
	"github.com/bro-n-bro/spacebox/broker/model"
)

func (m *Module) HandleBlock(ctx context.Context, block *types.Block) error {
	if err := m.broker.PublishBlock(ctx, model.Block{
		Height:          block.Height,
		Hash:            block.Hash,
		ProposerAddress: block.ProposerAddress,
		NumTxs:          int64(block.TxNum),
		TotalGas:        block.TotalGas,
		Timestamp:       block.Timestamp,
	}); err != nil {
		m.log.Error().
			Err(err).
			Int64("block_height", block.Height).
			Msg("PublishBlock error")

		return err
	}

	for _, precommit := range block.ValidatorPrecommits {
		if err := m.broker.PublishValidatorPrecommit(ctx, model.ValidatorPrecommit{
			Height:           block.Height,
			ValidatorAddress: precommit.ValidatorAddress,
			BlockIDFlag:      precommit.BlockIDFlag,
			Timestamp:        precommit.Timestamp,
		}); err != nil {
			return errors.Wrap(err, "PublishValidatorPrecommit error")
		}
	}

	return nil
}
