package worker

import (
	"context"
	"fmt"
	"sync"
	"syscall"
	"time"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/rs/zerolog"

	"github.com/bro-n-bro/spacebox-crawler/internal/rep"
	tb "github.com/bro-n-bro/spacebox-crawler/pkg/mapper/to_broker"
	ts "github.com/bro-n-bro/spacebox-crawler/pkg/mapper/to_storage"
	"github.com/bro-n-bro/spacebox-crawler/types"
)

type (
	Worker struct {
		log *zerolog.Logger
		wg  *sync.WaitGroup

		tsM        ts.ToStorage
		storage    rep.Storage
		cdc        codec.Codec
		tbM        tb.ToBroker
		broker     rep.Broker
		rpcClient  rep.RPCClient
		grpcClient rep.GrpcClient

		metrics *metrics

		stopProcessing         func()
		stopWsListener         func()
		stopEnqueueHeight      func()
		stopEnqueueErrorBlocks func()

		heightCh chan int64

		modules []types.Module
		cfg     Config
	}
	metrics struct {
		durMetric *prometheus.HistogramVec
	}
)

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
	if w.cfg.MetricsEnabled {
		w.metrics = &metrics{
			durMetric: promauto.NewHistogramVec(prometheus.HistogramOpts{
				Namespace: "spacebox_crawler",
				Name:      "process_duration",
				Help:      "Duration of parsed blockchain objects",
			}, []string{"type"}),
		}
	}

	ctx, cancel := context.WithCancel(context.Background())
	w.stopProcessing = cancel

	// workers count must be greater than 0
	workersCount := w.cfg.WorkersCount
	if workersCount <= 0 {
		workersCount = 1
	}

	w.heightCh = make(chan int64, workersCount)

	stopHeight := w.cfg.StopHeight
	// check if stop height is empty, and we want to process height range from config
	if stopHeight <= 0 && w.cfg.StartHeight >= 0 {
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
			return fmt.Errorf("failed to subscribe to new blocks: %w", err)
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
		if err := syscall.Kill(syscall.Getpid(), syscall.SIGINT); err != nil {
			panic(err)
		}
	}(wg)

	return nil
}

func (w *Worker) Stop(_ context.Context) error {
	w.stopEnqueueHeight()
	w.stopEnqueueErrorBlocks()

	if w.cfg.ProcessNewBlocks && w.rpcClient.WsEnabled() {
		w.stopWsListener()
	}

	t := time.NewTicker(2 * time.Second)
	defer t.Stop()
	<-t.C // XXX save from send to closed channel

	close(w.heightCh)
	w.wg.Wait()
	w.stopProcessing()

	w.log.Info().Msg("stop workers")

	return nil
}
