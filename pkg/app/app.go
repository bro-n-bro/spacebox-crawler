package app

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

type App interface {
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
}

// Run starts an application as the graceful shutdown service
func Run(a App) error {
	// start an application
	if err := a.Start(context.Background()); err != nil {
		return err
	}

	// wait for OS signal for graceful shutdown
	quitCh := make(chan os.Signal, 1)
	signal.Notify(quitCh, os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	<-quitCh

	// stop an application
	if err := a.Stop(context.Background()); err != nil {
		return err
	}

	return nil
}
