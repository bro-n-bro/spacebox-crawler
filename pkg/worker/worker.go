package worker

import (
	"context"
	"fmt"
	"os"
	"sync"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/rs/zerolog"

	"bro-n-bro-osmosis/internal/rep"
	tb "bro-n-bro-osmosis/pkg/mapper/to_broker"
	ts "bro-n-bro-osmosis/pkg/mapper/to_storage"
	"bro-n-bro-osmosis/types"
)

type Worker struct {
	log *zerolog.Logger
	wg  *sync.WaitGroup

	broker     rep.Broker
	rpcClient  rep.RPCClient
	grpcClient rep.GrpcClient
	storage    rep.Storage
	cdc        codec.Codec
	tbM        tb.ToBroker
	tsM        ts.ToStorage

	cfg Config

	modules []types.Module

	cancel   func()
	heightCh chan int64
}

func New(cfg Config, b rep.Broker, rpcCli rep.RPCClient, grpcCli rep.GrpcClient, modules []types.Module, s rep.Storage,
	marshaler codec.Codec, tbM tb.ToBroker, tsM ts.ToStorage) *Worker {

	l := zerolog.New(os.Stderr).Output(zerolog.ConsoleWriter{Out: os.Stderr}).With().Timestamp().
		Str("cmp", "worker").Logger()

	chanSize := cfg.ChanSize
	if chanSize == 0 {
		chanSize = 1
	}

	w := &Worker{
		cfg:        cfg,
		log:        &l,
		broker:     b,
		rpcClient:  rpcCli,
		grpcClient: grpcCli,
		storage:    s,
		modules:    modules,
		cdc:        marshaler,
		tbM:        tbM,
		tsM:        tsM,
		wg:         &sync.WaitGroup{},
		heightCh:   make(chan int64, chanSize),
	}
	w.fillModules()
	return w
}

func (w *Worker) Start(_ context.Context) error {
	ctx, cancel := context.WithCancel(context.Background())
	w.cancel = cancel

	workersCount := w.cfg.WorkersCount
	if workersCount == 0 {
		workersCount = 1
	}

	for i := 0; i < workersCount; i++ {
		w.wg.Add(1)
		go w.processHeight(ctx, i) // run processing block function
	}

	if w.cfg.ProcessNewBlocks && w.rpcClient.WsEnabled() {
		eventCh, err := w.rpcClient.SubscribeNewBlocks(ctx)
		if err != nil {
			return fmt.Errorf("failed to subscribe to new blocks: %s", err)
		}
		go w.enqueueNewBlocks(ctx, eventCh)
	}

	go w.enqueueHeight(ctx)

	return nil
}

func (w *Worker) Stop(_ context.Context) error {
	w.cancel()
	w.wg.Wait()
	close(w.heightCh)
	w.log.Info().Msg("stop workers")
	return nil
}
