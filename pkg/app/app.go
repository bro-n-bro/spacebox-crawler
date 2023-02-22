package app

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type App interface {
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
	GetStartTimeout() time.Duration
	GetStopTimeout() time.Duration
}

// Run starts an application as the graceful shutdown service
func Run(a App) error {
	startCtx, startCancel := context.WithTimeout(context.Background(), a.GetStartTimeout())
	defer startCancel()

	// start an application
	if err := a.Start(startCtx); err != nil {
		return err
	}

	// wait for OS signal for graceful shutdown
	quitCh := make(chan os.Signal, 1)
	signal.Notify(quitCh, os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	<-quitCh

	stopCtx, stopCancel := context.WithTimeout(context.Background(), a.GetStopTimeout())
	defer stopCancel()

	// stop an application
	if err := a.Stop(stopCtx); err != nil {
		return err
	}

	return nil
}
