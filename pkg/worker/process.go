package worker

import (
	"context"
	"encoding/json"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	tmtcoreypes "github.com/tendermint/tendermint/rpc/core/types"
	tmtypes "github.com/tendermint/tendermint/types"
	"golang.org/x/sync/errgroup"

	"bro-n-bro-osmosis/types"
)

func (w *Worker) processHeight(ctx context.Context, workerIndex int) {
	var parsedCount int
	defer w.wg.Done()
	//defer func() {
	//	log.Printf("worker: %d. parsed %d blocks", wNumber, parsedCount)
	//}()

	for {
		select {
		case <-ctx.Done():
			w.log.Info().Int("worker_number", workerIndex).Msg("done worker")
			return
		case height := <-w.heightCh:
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
			}

			w.log.Info().Int("worker_number", workerIndex).Msgf("Parse block № %d", height)

			g, _ctx := errgroup.WithContext(ctx)

			var (
				block *tmtcoreypes.ResultBlock
				vals  *tmtcoreypes.ResultValidators
			)

			g.Go(func() error {
				var err error

				_blockDur := time.Now()
				block, err = w.grpcClient.Block(_ctx, height)
				if err != nil {
					return err
				}
				w.log.Info().
					Int("worker_number", workerIndex).
					Int64("block_height", height).
					Dur("parse_block_dur", time.Since(_blockDur)).
					Msg("Get block info")
				return nil
			})

			g.Go(func() error {
				var err error

				_validatorsDur := time.Now()
				vals, err = w.grpcClient.Validators(_ctx, height)
				if err != nil {
					return err
				}
				w.log.Info().
					Int("worker_number", workerIndex).
					Int64("block_height", height).
					Dur("parse_block_dur", time.Since(_validatorsDur)).
					Msg("Get validators info")
				return nil
			})

			if err := g.Wait(); err != nil {
				w.log.Error().Err(err).Msgf("processHeight block got error: %v", err)
				continue
			}

			_txsDur := time.Now()
			// with async
			// 2022/11/03 17:29:29 worker 7. height: 6710901. get block dur: 854.136709ms
			// 2022/11/03 17:29:35 worker 7. height: 6710901. get txs dur: 5.954485916s
			// 2022/11/03 17:29:36 worker 7. height: 6710901. get validators dur: 536.417667ms

			// sync
			// 2022/11/03 17:30:02 worker 7. height: 6710901. get block dur: 645.112792ms
			// 2022/11/03 17:30:08 worker 7. height: 6710901. get txs dur: 5.895915416s
			// 2022/11/03 17:30:08 worker 7. height: 6710901. get validators dur: 434.883833ms

			// Async: 2022/11/03 19:11:21 duration sec: 208.474833083
			// Sync: 2022/11/03 19:14:38 duration sec: 165.029986125

			txsRes, err := w.grpcClient.TxsOld(ctx, block.Block.Data.Txs)
			if err != nil {
				w.log.Error().Err(err).Msgf("get txs error: %v", err)
				continue
			}
			w.log.Info().
				Int("worker_number", workerIndex).
				Int64("block_height", height).
				Dur("txs_dur", time.Since(_txsDur)).
				Msg("Get txs info")

			b := types.NewBlockFromTmBlock(block)
			txs := types.NewTxsFromTmTxs(txsRes)

			if err := w.broker.PublishBlock(ctx, w.tbM.MapBlock(b, txs.TotalGas())); err != nil {
				w.log.Error().
					Err(err).
					Int64("block_height", height).
					Msg("PublishBlock error")
				continue
			}

			// handle block first
			w.processBlock(ctx, b, vals)
			w.processTxs(ctx, txs)

			_ = vals
			_ = txs
			parsedCount++
		}
	}
}

func (w *Worker) processGenesis(ctx context.Context, genesis *tmtypes.GenesisDoc) error {
	var appState map[string]json.RawMessage
	if err := json.Unmarshal(genesis.AppState, &appState); err != nil {
		w.log.Err(err).Msgf("error unmarshalling genesis doc: %v", err)
		return err
	}

	for _, module := range w.modules {
		if genesisModule, ok := module.(types.GenesisModule); ok {
			if err := genesisModule.HandleGenesis(ctx, genesis, appState); err != nil {
				w.log.Error().Err(err).Msgf("handle genesis error. module: %s, err: %v", module, err)
			}
		}
	}
	return nil
}

func (w *Worker) processBlock(ctx context.Context, block *types.Block, vals *tmtcoreypes.ResultValidators) {
	for _, m := range w.modules {
		hbI, ok := m.(types.BlockModule)
		if ok {
			if err := hbI.HandleBlock(ctx, block, vals); err != nil {
				w.log.Error().Err(err).Msgf("HandleBlock error:", err)
			}
		}
	}
}

func (w *Worker) processTxs(ctx context.Context, txs []*types.Tx) {
	for _, tx := range txs {
		for _, m := range w.modules {
			hTxI, ok := m.(types.TransactionModule)
			if ok {
				if err := hTxI.HandleTx(ctx, tx); err != nil {
					w.log.Error().Err(err).Msgf("HandleTX error:", err)
					continue
				}
			}
		}

		for i, msg := range tx.Body.Messages {
			var stdMsg sdk.Msg

			err := w.marshaler.UnpackAny(msg, &stdMsg)
			if err != nil {
				w.log.Error().Err(err).Msgf("error while unpacking message: %s", err)
				continue
				//return fmt.Errorf("error while unpacking message: %s", err)
			}

			for _, m := range w.modules {
				if messageModule, ok := m.(types.MessageModule); ok {
					err = messageModule.HandleMessage(ctx, i, stdMsg, tx)
					if err != nil {
						w.log.Error().Err(err).Msgf("HandleMessage error:", m, tx, stdMsg, err)
						continue
					}
				}
			}
		}
	}
}
