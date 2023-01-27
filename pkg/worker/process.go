package worker

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	tmtcoreypes "github.com/tendermint/tendermint/rpc/core/types"
	tmtypes "github.com/tendermint/tendermint/types"
	"golang.org/x/sync/errgroup"

	"github.com/bro-n-bro/spacebox-crawler/types"
)

const (
	keyModule = "module"
)

func (w *Worker) processHeight(ctx context.Context, workerIndex int) {
	var parsedCount int
	defer w.wg.Done()
	defer func() {
		w.log.Debug().Msgf("worker: %d. parsed %d blocks", workerIndex, parsedCount)
	}()

	for height := range w.heightCh {
		select {
		case <-ctx.Done():
			w.log.Info().Int("worker_index", workerIndex).Msg("done worker")
			return
		default:
		}

		// for debug
		parsedCount++

		if height == 0 {
			// TODO: in bdjuno genesis can be gets from local path
			w.log.Info().Int("worker_number", workerIndex).Msg("Parse genesis")

			_genesisDur := time.Now()

			genesis, err := w.rpcClient.Genesis(ctx)
			if err != nil {
				w.log.Error().Err(err).Msgf("get genesis error: %v", err)
				continue
			}

			w.log.Info().
				Int("worker_number", workerIndex).
				Msgf("Get genesis dur: %v", time.Since(_genesisDur))

			if err = w.processGenesis(ctx, genesis); err != nil {
				w.log.Fatal().Err(err).Msgf("processHeight genesis error %v:", err)
				continue
			}

			continue
		} // TODO: what about storage?

		if err := w.checkOrCreateBlockInStorage(ctx, height); err != nil {
			switch {
			case errors.Is(err, ErrBlockProcessed):
				w.log.Debug().Int64("height", height).Msg("block already processed. skip height")
			case errors.Is(err, ErrBlockProcessing):
				w.log.Debug().Int64("height", height).Msg("block is already processing now. skip height")
			case errors.Is(err, ErrBlockError):
				w.log.Debug().Int64("height", height).Msg("block processed with error. " +
					"if you want to process this height again see PROCESS_ERROR_BLOCKS ENV")
			}

			continue
		}

		w.log.Info().Int("worker_number", workerIndex).Msgf("Parse block № %d", height)

		g, ctx2 := errgroup.WithContext(ctx)

		var (
			block *tmtcoreypes.ResultBlock
			vals  *tmtcoreypes.ResultValidators
		)

		g.Go(func() error {
			var err error

			_blockDur := time.Now()
			if block, err = w.grpcClient.Block(ctx2, height); err != nil {
				return err
			}
			w.log.Debug().
				Int("worker_number", workerIndex).
				Int64("block_height", height).
				Dur("parse_block_dur", time.Since(_blockDur)).
				Msg("Get block info")
			return nil
		})

		g.Go(func() error {
			var err error
			_validatorsDur := time.Now()
			if vals, err = w.grpcClient.Validators(ctx2, height); err != nil {
				return err
			}
			w.log.Debug().
				Int("worker_number", workerIndex).
				Int64("block_height", height).
				Dur("parse_block_dur", time.Since(_validatorsDur)).
				Msg("Get validators info")

			return nil
		})

		if err := g.Wait(); err != nil {
			w.log.Error().Err(err).Msgf("processHeight block got error: %v", err)
			w.setErrorStatusWithLogging(ctx, height, err.Error())

			continue
		}

		_txsDur := time.Now()

		txsRes, err := w.grpcClient.Txs(ctx, block.Block.Data.Txs)
		if err != nil {
			w.log.Error().Err(err).Msgf("get txs error: %v", err)
			w.setErrorStatusWithLogging(ctx, height, err.Error())

			continue
		}

		w.log.Debug().
			Int("worker_number", workerIndex).
			Int64("block_height", height).
			Dur("txs_dur", time.Since(_txsDur)).
			Msg("Get txs info")

		txs := types.NewTxsFromTmTxs(txsRes, w.cdc)
		b := types.NewBlockFromTmBlock(block, txs.TotalGas())

		g, ctx2 = errgroup.WithContext(ctx)

		g.Go(func() error {
			return w.withMetrics("validators", func() error {
				return w.processValidators(ctx2, vals)
			})
		})
		g.Go(func() error {
			return w.withMetrics("block", func() error {
				return w.processBlock(ctx2, b)
			})
		})
		g.Go(func() error {
			return w.withMetrics("txs", func() error {
				return w.processTxs(ctx2, txs)
			})
		})
		g.Go(func() error {
			return w.withMetrics("messages", func() error {
				return w.processMessages(ctx2, txs)
			})
		})

		if err := g.Wait(); err != nil {
			w.setErrorStatusWithLogging(ctx, height, err.Error())
			continue
		}

		if err := w.storage.SetProcessedStatus(ctx, height); err != nil {
			w.log.Error().Err(err).Int64("height", height).Msgf("cant set processed status in storage %v:", err)
		}
	}
}

func (w *Worker) processGenesis(ctx context.Context, genesis *tmtypes.GenesisDoc) error {
	var appState map[string]json.RawMessage
	if err := json.Unmarshal(genesis.AppState, &appState); err != nil {
		w.log.Err(err).Msgf("error unmarshalling genesis doc: %v", err)
		return err
	}

	for _, m := range genesisHandlers {
		if err := m.HandleGenesis(ctx, genesis, appState); err != nil {
			w.log.Error().Err(err).Str(keyModule, m.Name()).Msgf("handle genesis error: %v", err)
		}
	}

	return nil
}

func (w *Worker) processBlock(ctx context.Context, block *types.Block) error {
	for _, m := range blockHandlers {
		if err := m.HandleBlock(ctx, block); err != nil {
			w.log.Error().Str(keyModule, m.Name()).Err(err).Msgf("HandleBlock error: %v", err)
			return err
		}
	}

	return nil
}

func (w *Worker) processValidators(ctx context.Context, vals *tmtcoreypes.ResultValidators) error {
	for _, m := range validatorsHandlers {
		if err := m.HandleValidators(ctx, vals); err != nil {
			w.log.Error().Str(keyModule, m.Name()).Err(err).Msgf("HandleValidators error: %v", err)
			return err
		}
	}

	return nil
}

func (w *Worker) processTxs(ctx context.Context, txs []*types.Tx) error {
	for _, tx := range txs {
		for _, m := range transactionHandlers {
			if err := m.HandleTx(ctx, tx); err != nil {
				w.log.Error().Str(keyModule, m.Name()).Err(err).Msgf("HandleTX error %v", err)
				return err
			}
		}
	}

	return nil
}

func (w *Worker) processMessages(ctx context.Context, txs []*types.Tx) error {
	for _, tx := range txs {
		if !tx.Successful() { // skip message processing from failed transaction
			continue
		}

		for i, msg := range tx.Body.Messages {
			var stdMsg sdk.Msg

			if err := w.cdc.UnpackAny(msg, &stdMsg); err != nil {
				w.log.Error().Err(err).Msgf("error while unpacking message: %s", err)
				return err
			}

			for _, m := range messageHandlers {
				if err := m.HandleMessage(ctx, i, stdMsg, tx); err != nil {
					w.log.Error().Str(keyModule, m.Name()).Err(err).Msgf("HandleMessage error: %v", err)
					return err
				}
			}
		}
	}

	return nil
}
