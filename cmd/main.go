package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/caarlos0/env/v6"
	"github.com/cosmos/cosmos-sdk/simapp"
	"github.com/cosmos/cosmos-sdk/simapp/params"
	"github.com/cosmos/cosmos-sdk/std"
	sdk "github.com/cosmos/cosmos-sdk/types"
	cosmomodule "github.com/cosmos/cosmos-sdk/types/module"
	"github.com/joho/godotenv"
	tmtcoreypes "github.com/tendermint/tendermint/rpc/core/types"
	tmtypes "github.com/tendermint/tendermint/types"
	"golang.org/x/sync/errgroup"

	"bro-n-bro-osmosis/adapter/broker"
	grpcClient "bro-n-bro-osmosis/client/grpc"
	"bro-n-bro-osmosis/client/rpc"
	"bro-n-bro-osmosis/internal/app"
	"bro-n-bro-osmosis/modules"
	"bro-n-bro-osmosis/modules/messages"
	"bro-n-bro-osmosis/types"
)

var (
	encoding = params.EncodingConfig{}
)

// getBasicManagers returns the various basic managers that are used to register the encoding to
// support custom messages.
// This should be edited by custom implementations if needed.
func getBasicManagers() []cosmomodule.BasicManager {
	return []cosmomodule.BasicManager{
		simapp.ModuleBasics,
	}
}

// MakeEncodingConfig creates an EncodingConfig to properly handle all the messages
func MakeEncodingConfig(managers []cosmomodule.BasicManager) params.EncodingConfig {

	encodingConfig := params.MakeTestEncodingConfig()
	std.RegisterLegacyAminoCodec(encodingConfig.Amino)
	std.RegisterInterfaces(encodingConfig.InterfaceRegistry)
	manager := mergeBasicManagers(managers)
	manager.RegisterLegacyAminoCodec(encodingConfig.Amino)
	manager.RegisterInterfaces(encodingConfig.InterfaceRegistry)
	return encodingConfig
}

// mergeBasicManagers merges the given managers into a single module.BasicManager
func mergeBasicManagers(managers []cosmomodule.BasicManager) cosmomodule.BasicManager {
	var union = cosmomodule.BasicManager{}
	for _, manager := range managers {
		for k, v := range manager {
			union[k] = v
		}
	}
	return union
}

func main() {
	//oldMain()
	newMain()
}

func newMain() {
	// load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		panic(err)
	}

	var cfg app.Config
	// fill these variables into a struct
	if err := env.Parse(&cfg); err != nil {
		panic(err)
	}

	a := app.New(cfg)
	if err := a.Start(context.Background()); err != nil {
		panic(err)
	}

	quitCh := make(chan os.Signal, 1)
	signal.Notify(quitCh, os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	<-quitCh

	if err := a.Stop(context.Background()); err != nil {
		panic(err)
	}

	fmt.Printf("%+v\n", cfg)
}

func oldMain() {
	encoding = MakeEncodingConfig(getBasicManagers())

	grpcCli := grpcClient.New("195.201.56.108:9090")
	//grpcCli := grpcClient.New("grpc.osmo-test-4.cybernode.ai:1443")
	//grpcCli := grpcClient.New("157.90.93.137:9090")
	//grpcCli := grpcClient.New("95.216.241.52:26090")
	//grpcCli := grpcClient.New("cosmos-grpc.polkachu.com:14990")
	//grpcCli := grpcClient.New("grpc-test.osmosis.zone:443")

	ctx := context.Background()
	if err := grpcCli.Start(ctx); err != nil {
		log.Fatal(err)
	}
	defer grpcCli.Stop(ctx)

	rpcCli := rpc.New("https://rpc.bostrom.bronbro.io:443", false)
	if err := rpcCli.Start(ctx); err != nil {
		log.Fatal(err)
	}
	defer rpcCli.Stop(ctx)

	b := broker.New()
	if err := b.Start(ctx); err != nil {
		log.Fatal(err)
	}
	defer b.Stop(ctx)

	parser := messages.JoinMessageParsers(messages.CosmosMessageAddressesParser)

	curModules := modules.BuildModules(b, grpcCli, parser, encoding.Marshaler, "bank", "gov", "auth", "mint",
		"slashing", "staking")

	ch := make(chan int64, 8)
	wg := &sync.WaitGroup{}
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go worker(wg, i, grpcCli, rpcCli, curModules, ch)
	}

	start := time.Now()

	ch <- 0
	for i := int64(100); i < 200; i++ {
		ch <- i
	}

	// osmosis min
	//ch <- 6286335

	close(ch)
	wg.Wait()

	log.Println("duration sec:", time.Since(start).Seconds())

}

type toCrowler struct {
	block *types.Block
	txs   []*types.Tx
	val   *types.Validators
}

func worker(wg *sync.WaitGroup, wNumber int, cli *grpcClient.Client, rpcCli *rpc.Client,
	modules []types.Module, ch chan int64) {

	defer wg.Done()
	var parsedCount int
	//defer func() {
	//	log.Printf("worker: %d. parsed %d blocks", wNumber, parsedCount)
	//}()

	for height := range ch {

		if height == 0 {
			// FIXME: in bdjuno genesis can be gets from local path
			log.Printf("worker: %d. Parse genesis", wNumber)

			_genesisDur := time.Now()
			genesis, err := rpcCli.RpcClient.Genesis(context.Background())
			if err != nil {
				log.Fatal("get genesis error: ", err)
			}
			log.Printf("worker %d. get genesis dur: %v", wNumber, time.Since(_genesisDur))

			if err = processGenesis(context.Background(), modules, genesis.Genesis); err != nil {
				log.Fatal("process genesis error:", err)
				return
			}

			continue
		}

		log.Printf("worker: %d. Parse block № %d", wNumber, height)

		g, _ctx := errgroup.WithContext(context.Background())

		var (
			block *tmtcoreypes.ResultBlock
			vals  *tmtcoreypes.ResultValidators
		)

		g.Go(func() error {
			var err error

			_blockDur := time.Now()
			block, err = cli.Block(_ctx, height)
			if err != nil {
				return err
			}
			log.Printf("worker %d. height: %v. get block dur: %v", wNumber, height, time.Since(_blockDur))
			return nil
		})

		g.Go(func() error {
			var err error

			_validatorsDur := time.Now()
			vals, err = cli.Validators(_ctx, height)
			if err != nil {
				return err
			}
			log.Printf("worker %d. height: %v. get validators dur: %v", wNumber, height, time.Since(_validatorsDur))
			return nil
		})

		if err := g.Wait(); err != nil {
			log.Fatal(err)
			return
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
		txs, err := cli.TxsOld(context.Background(), block.Block.Data.Txs)
		if err != nil {
			log.Fatal(err)
			return
		}
		log.Printf("worker %d. height: %v. get txs dur: %v", wNumber, height, time.Since(_txsDur))

		// handle block first
		processBlock(modules, types.NewBlockFromTmBlock(block), vals)
		processTxs(modules, types.NewTxsFromTmTxs(txs))

		_ = vals
		_ = toCrowler{
			//block: resp,
			//txs: txs,
			//val: vals,
		}
		parsedCount++
	}
}

func processGenesis(ctx context.Context, mds []types.Module, genesis *tmtypes.GenesisDoc) error {
	var appState map[string]json.RawMessage
	if err := json.Unmarshal(genesis.AppState, &appState); err != nil {
		log.Fatalf("error unmarshalling genesis doc: %s", err)
		return err
	}

	for _, module := range mds {
		if genesisModule, ok := module.(types.GenesisModule); ok {
			if err := genesisModule.HandleGenesis(ctx, genesis, appState); err != nil {
				log.Printf("handle genesis error. module: %s, err: %v", module, err)
			}
		}
	}
	return nil
}

func processBlock(mds []types.Module, block *types.Block, vals *tmtcoreypes.ResultValidators) {
	for _, m := range mds {
		hbI, ok := m.(types.BlockModule)
		if ok {
			if err := hbI.HandleBlock(context.Background(), block, vals); err != nil {
				log.Println("HandleBlock error:", err)
			}
		}
	}
}

func processTxs(mds []types.Module, txs []*types.Tx) {
	for _, tx := range txs {
		for _, m := range mds {
			hTxI, ok := m.(types.TransactionModule)
			if ok {
				if err := hTxI.HandleTx(context.Background(), tx); err != nil {
					log.Println("HandleTX error:", err)
					continue
				}
			}
		}

		for i, msg := range tx.Body.Messages {
			var stdMsg sdk.Msg

			err := encoding.Marshaler.UnpackAny(msg, &stdMsg)
			if err != nil {
				log.Printf("error while unpacking message: %s", err)
				continue
				//return fmt.Errorf("error while unpacking message: %s", err)
			}

			for _, m := range mds {
				if messageModule, ok := m.(types.MessageModule); ok {
					err = messageModule.HandleMessage(context.Background(), i, stdMsg, tx)
					if err != nil {
						log.Println(m, tx, stdMsg, err)
					}
				}
			}
		}
	}
}
