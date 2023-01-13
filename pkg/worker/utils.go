package worker

import (
	"context"

	"github.com/pkg/errors"
)

var (
	ErrBlockProcessed  = errors.New("block already processed")
	ErrBlockProcessing = errors.New("block is processing right now")
	ErrBlockError      = errors.New("block processed with error")
)

func (w *Worker) setErrorStatusWithLogging(ctx context.Context, height int64, msg string) {
	if err := w.storage.SetErrorStatus(ctx, height, msg); err != nil {
		w.log.Error().Err(err).Int64("height", height).Msgf("cant set error status in storage %v:", err)
	}
}

func (w *Worker) checkOrCreateBlockInStorage(ctx context.Context, height int64) error {
	hasBlock, err := w.storage.HasBlock(ctx, height)
	if err != nil {
		w.log.Fatal().Err(err).Int64("height", height).Msgf("cant check block in storage %v:", err)
		return err
	}

	if hasBlock { // exists block
		status, err := w.storage.GetBlockStatus(ctx, height)
		if err != nil {
			w.log.Error().Err(err).Int64("height", height).Msgf("cant get block status in storage %v:", err)
			return err
		}

		switch {
		// block info already in kafka
		case status.IsProcessed():
			return ErrBlockProcessed
		// block now is processing
		case status.IsProcessing():
			return ErrBlockProcessing
		// block processed with error, skip if needed
		case status.IsError() && !w.cfg.ProcessErrorBlocks:
			return ErrBlockError
		}
	} else { // new block
		if err := w.storage.CreateBlock(ctx, w.tsM.NewBlock(height)); err != nil {
			w.log.Error().Err(err).Int64("height", height).Msgf("cant create new block in storage %v:", err)
			return err
		}
	}

	return nil
}
