package worker

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/pkg/errors"

	"github.com/cosmos/cosmos-sdk/codec"
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

	cancel   func()
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
	w.cancel = cancel

	if err := w.pingStorage(ctx); err != nil {
		return err
	}

	workersCount := w.cfg.WorkersCount
	if workersCount == 0 {
		workersCount = 1
	}

	w.heightCh = make(chan int64, workersCount)

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
