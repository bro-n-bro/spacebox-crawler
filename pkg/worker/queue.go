package worker

import (
	"context"

	tmtcoreypes "github.com/tendermint/tendermint/rpc/core/types"
	tmtypes "github.com/tendermint/tendermint/types"
)

func (w *Worker) enqueueHeight(ctx context.Context) {
	for height := w.cfg.StartHeight; height <= w.cfg.StopHeight; height++ {
		// safe from closed channel
		select {
		case <-ctx.Done():
			return
		default:
		}
		w.heightCh <- height // put height to channel for processing the block
	}

	if w.cfg.ProcessErrorBlocks {
		heights, err := w.storage.GetErrorBlockHeights(ctx)
		if err != nil {
			w.log.Error().Err(err).Msg("GetErrorBlockHeights error")
		} else {
			for _, height := range heights {
				// safe from closed channel
				select {
				case <-ctx.Done():
					return
				default:
				}
				w.heightCh <- height
			}
		}
	}

	if !w.cfg.ProcessNewBlocks || !w.rpcClient.WsEnabled() {
		// TODO: stop program
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
			w.log.Info().Int64("height", height).Msgf("enqueueing new block with height: %v", height)
			w.heightCh <- height
		}
	}
}
