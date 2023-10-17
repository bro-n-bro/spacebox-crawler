package worker

import (
	"context"
	"strings"

	codec "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/pkg/errors"

	"github.com/bro-n-bro/spacebox-crawler/adapter/storage/model"
	"github.com/bro-n-bro/spacebox-crawler/types"
)

var (
	ErrBlockProcessed  = errors.New("block already processed")
	ErrBlockProcessing = errors.New("block is processing right now")
	ErrBlockError      = errors.New("block processed with error")
)

func (w *Worker) setErrorStatusWithLogging(ctx context.Context, height int64, msg string) {
	if err := w.storage.SetErrorStatus(ctx, height, msg); err != nil {
		w.log.Error().Err(err).Int64("height", height).Msg("can't set error status in storage")
	}
}

func (w *Worker) checkOrCreateBlockInStorage(ctx context.Context, height int64) error {
	block, err := w.storage.GetBlockByHeight(ctx, height)
	if err != nil && errors.Is(err, types.ErrBlockNotFound) {
		// create new block
		if err = w.storage.CreateBlock(ctx, w.tsM.NewBlock(height)); err != nil {
			w.log.Error().Err(err).Int64("height", height).Msg("can't create new block in storage")
			return err
		}

		return nil
	} else if err != nil {
		// got some error from storage
		return err
	}

	// block exists check status
	switch {
	// block info already in kafka
	case block.Status.IsProcessed():
		return ErrBlockProcessed
	// block now is processing
	case block.Status.IsProcessing():
		return ErrBlockProcessing
	// block processed with error, skip if needed
	case block.Status.IsError():
		if !w.cfg.ProcessErrorBlocks {
			return ErrBlockError
		}
		return w.storage.UpdateStatus(ctx, height, model.StatusProcessing)
	}
	return nil
}

func (w *Worker) unpackMessage(ctx context.Context, height int64, msg *codec.Any) (stdMsg sdk.Msg, err error) {
	if err = w.cdc.UnpackAny(msg, &stdMsg); err == nil {
		return stdMsg, nil
	}

	if strings.HasPrefix(err.Error(), "no concrete type registered for type URL") {
		w.log.Warn().Err(err).Msg("error while unpacking message")

		if err = w.storage.InsertErrorMessage(ctx, w.tsM.NewErrorMessage(height, err.Error())); err != nil {
			w.log.Error().
				Err(err).
				Int64(keyHeight, height).
				Msg("fail to insert error_message")
			return nil, err
		}

		// just skip unsupported message
		return nil, nil
	}

	w.log.Error().Err(err).Msg("error while unpacking message")

	return nil, err
}
