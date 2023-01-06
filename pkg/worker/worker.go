package worker

import (
	"context"
	"fmt"
	"sync"
	"syscall"
	"time"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	"github.com/hexy-dev/spacebox-crawler/internal/rep"
	tb "github.com/hexy-dev/spacebox-crawler/pkg/mapper/to_broker"
	ts "github.com/hexy-dev/spacebox-crawler/pkg/mapper/to_storage"
	"github.com/hexy-dev/spacebox-crawler/types"
)

type Worker struct {
	log *zerolog.Logger
	wg  *sync.WaitGroup

	tsM        ts.ToStorage
	storage    rep.Storage
	cdc        codec.Codec
	tbM        tb.ToBroker
	broker     rep.Broker
	rpcClient  rep.RPCClient
	grpcClient rep.GrpcClient

	stopProcessing         func()
	stopWsListener         func()
	stopEnqueueHeight      func()
	stopEnqueueErrorBlocks func()

	heightCh chan int64

	modules []types.Module
	cfg     Config
}

func New(cfg Config, l zerolog.Logger, b rep.Broker, rpcCli rep.RPCClient, grpcCli rep.GrpcClient,
	modules []types.Module, s rep.Storage, marshaler codec.Codec, tbM tb.ToBroker, tsM ts.ToStorage) *Worker {
	l = l.With().Str("cmp", "worker").Logger()

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
	}
	w.fillModules()
	return w
}

func (w *Worker) Start(_ context.Context) error {
	ctx, cancel := context.WithCancel(context.Background())
	w.stopProcessing = cancel

	// workers count must be greater than 0
	workersCount := w.cfg.WorkersCount
	if workersCount == 0 {
		workersCount = 1
	}

	w.heightCh = make(chan int64, workersCount)

	// check if stop height is empty
	stopHeight := w.cfg.StopHeight
	if stopHeight == 0 {
		var err error
		stopHeight, err = w.rpcClient.GetLastBlockHeight(ctx)
		if err != nil {
			return err
		}
	}

	// spawn workers
	for i := 0; i < workersCount; i++ {
		w.wg.Add(1)
		go w.processHeight(ctx, i) // run processing block function
	}

	// subscribe to process new blocks by websocket
	if w.cfg.ProcessNewBlocks && w.rpcClient.WsEnabled() {
		eventCh, err := w.rpcClient.SubscribeNewBlocks(ctx)
		if err != nil {
			return fmt.Errorf("failed to subscribe to new blocks: %s", err)
		}
		go w.enqueueNewBlocks(ctx, eventCh)
	}

	wg := &sync.WaitGroup{}

	// enqueue error blocks height
	if w.cfg.ProcessErrorBlocks {
		wg.Add(1)
		go w.enqueueErrorBlocks(ctx, wg)
	}

	// enqueue block height based on config start/stop heights
	wg.Add(1)
	go w.enqueueHeight(ctx, wg, w.cfg.StartHeight, stopHeight)

	// graceful shutdown the application if processing is done
	go func(wg *sync.WaitGroup) {
		if w.cfg.ProcessNewBlocks && w.rpcClient.WsEnabled() { // we want to process new blocks
			w.log.Info().Msg("exit not needed")
			return
		}
		wg.Wait()
		w.log.Info().Msg("process block height done! stop program")
		_ = syscall.Kill(syscall.Getpid(), syscall.SIGINT)
	}(wg)

	return nil
}

func (w *Worker) Stop(_ context.Context) error {
	w.stopEnqueueHeight()
	w.stopEnqueueErrorBlocks()

	if w.cfg.ProcessNewBlocks && w.rpcClient.WsEnabled() {
		w.stopWsListener()
	}
	time.Sleep(1 * time.Second) // XXX save from send to closed channel
	close(w.heightCh)
	w.wg.Wait()
	w.stopProcessing()

	w.log.Info().Msg("stop workers")
	return nil
}

func (w *Worker) pingStorage(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, time.Minute)
	defer cancel()

	for {
		select {
		case <-ctx.Done():
			return errors.New("worker ping storage timeout")
		default:
			time.Sleep(1 * time.Second)
			if err := w.storage.Ping(ctx); err == nil {
				return nil
			}
		}
	}
}
