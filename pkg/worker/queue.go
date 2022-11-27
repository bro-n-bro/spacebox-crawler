package worker

import (
	"context"

	tmtcoreypes "github.com/tendermint/tendermint/rpc/core/types"
	tmtypes "github.com/tendermint/tendermint/types"
)

func (w *Worker) enqueueHeight() {
	for i := w.cfg.StartHeight; i < w.cfg.StopHeight; i++ {
		w.heightCh <- i // put height to channel for processing the block
	}
}

func (w *Worker) enqueueNewBlocks(ctx context.Context, eventCh <-chan tmtcoreypes.ResultEvent) {
	w.log.Info().Msg("listening for new block events...")

	for {
		select {
		case <-ctx.Done():
			w.log.Info().Msg("stop new block listener")
			return
		case e := <-eventCh:
			newBlock := e.Data.(tmtypes.EventDataNewBlock).Block
			height := newBlock.Header.Height
			w.log.Info().Int64("height", height).Msgf("enqueueing new block with height:", height)
			w.heightCh <- height
		}
	}
}
