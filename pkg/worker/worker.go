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
	"bro-n-bro-osmosis/types"
)

type Worker struct {
	log *zerolog.Logger
	wg  *sync.WaitGroup

	broker     rep.Broker
	rpcClient  rep.RPCClient
	grpcClient rep.GrpcClient
	cdc        codec.Codec
	tbM        tb.ToBroker

	cfg Config

	modules []types.Module

	cancel   func()
	heightCh chan int64
}

func New(cfg Config, b rep.Broker, rpcCli rep.RPCClient, grpcCli rep.GrpcClient, modules []types.Module,
	marshaler codec.Codec, tbM tb.ToBroker) *Worker {

	l := zerolog.New(os.Stderr).Output(zerolog.ConsoleWriter{Out: os.Stderr}).With().Timestamp().
		Str("cmp", "worker").Logger()
	w := &Worker{
		cfg:        cfg,
		log:        &l,
		broker:     b,
		rpcClient:  rpcCli,
		grpcClient: grpcCli,
		modules:    modules,
		cdc:        marshaler,
		tbM:        tbM,
		wg:         &sync.WaitGroup{},
		heightCh:   make(chan int64, cfg.ChanSize),
	}
	w.fillModules()
	return w
}

func (w *Worker) Start(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	w.cancel = cancel

	for i := 0; i < w.cfg.WorkersCount; i++ {
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

	go w.enqueueHeight()

	return nil

}

func (w *Worker) Stop(_ context.Context) error {
	w.cancel()
	w.wg.Wait()
	close(w.heightCh)
	w.log.Info().Msg("stop workers")
	return nil
}
