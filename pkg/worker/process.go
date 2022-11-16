package worker

//
//import (
//	"bro-n-bro-osmosis/types"
//	"context"
//	"log"
//	"time"
//
//	tmtcoreypes "github.com/tendermint/tendermint/rpc/core/types"
//	"golang.org/x/sync/errgroup"
//)
//
//func (w *Worker) process() {
//	defer wg.Done()
//	var parsedCount int
//	//defer func() {
//	//	log.Printf("worker: %d. parsed %d blocks", wNumber, parsedCount)
//	//}()
//
//	for height := range ch {
//
//		if height == 0 {
//			// FIXME: in bdjuno genesis can be gets from local path
//			log.Printf("worker: %d. Parse genesis", wNumber)
//
//			_genesisDur := time.Now()
//			genesis, err := rpcCli.RpcClient.Genesis(context.Background())
//			if err != nil {
//				log.Fatal("get genesis error: ", err)
//			}
//			log.Printf("worker %d. get genesis dur: %v", wNumber, time.Since(_genesisDur))
//
//			if err = processGenesis(context.Background(), modules, genesis.Genesis); err != nil {
//				log.Fatal("process genesis error:", err)
//				return
//			}
//
//			continue
//		}
//
//		log.Printf("worker: %d. Parse block № %d", wNumber, height)
//
//		g, _ctx := errgroup.WithContext(context.Background())
//
//		var (
//			block *tmtcoreypes.ResultBlock
//			vals  *tmtcoreypes.ResultValidators
//		)
//
//		g.Go(func() error {
//			var err error
//
//			_blockDur := time.Now()
//			block, err = cli.Block(_ctx, height)
//			if err != nil {
//				return err
//			}
//			log.Printf("worker %d. height: %v. get block dur: %v", wNumber, height, time.Since(_blockDur))
//			return nil
//		})
//
//		g.Go(func() error {
//			var err error
//
//			_validatorsDur := time.Now()
//			vals, err = cli.Validators(_ctx, height)
//			if err != nil {
//				return err
//			}
//			log.Printf("worker %d. height: %v. get validators dur: %v", wNumber, height, time.Since(_validatorsDur))
//			return nil
//		})
//
//		if err := g.Wait(); err != nil {
//			log.Fatal(err)
//			return
//		}
//
//		_txsDur := time.Now()
//		// with async
//		// 2022/11/03 17:29:29 worker 7. height: 6710901. get block dur: 854.136709ms
//		// 2022/11/03 17:29:35 worker 7. height: 6710901. get txs dur: 5.954485916s
//		// 2022/11/03 17:29:36 worker 7. height: 6710901. get validators dur: 536.417667ms
//
//		// sync
//		// 2022/11/03 17:30:02 worker 7. height: 6710901. get block dur: 645.112792ms
//		// 2022/11/03 17:30:08 worker 7. height: 6710901. get txs dur: 5.895915416s
//		// 2022/11/03 17:30:08 worker 7. height: 6710901. get validators dur: 434.883833ms
//
//		// Async: 2022/11/03 19:11:21 duration sec: 208.474833083
//		// Sync: 2022/11/03 19:14:38 duration sec: 165.029986125
//		txs, err := cli.TxsOld(context.Background(), block.Block.Data.Txs)
//		if err != nil {
//			log.Fatal(err)
//			return
//		}
//		log.Printf("worker %d. height: %v. get txs dur: %v", wNumber, height, time.Since(_txsDur))
//
//		// handle block first
//		processBlock(modules, types.NewBlockFromTmBlock(block), vals)
//		processTxs(modules, types.NewTxsFromTmTxs(txs))
//
//		_ = vals
//		_ = toCrowler{
//			//block: resp,
//			//txs: txs,
//			//val: vals,
//		}
//		parsedCount++
//	}
//}
