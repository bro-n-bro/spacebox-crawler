package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"

	"bro-n-bro-osmosis/internal/app"
)

func main() {
	_main()
}

func _main() {
	// load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		panic(err)
	}

	var cfg app.Config
	// fill these variables into a struct
	if err := env.Parse(&cfg); err != nil {
		panic(err)
	}

	// create an application
	a := app.New(cfg)

	// start an application
	if err := a.Start(context.Background()); err != nil {
		panic(err)
	}

	// wait for OS signal for graceful shutdown
	quitCh := make(chan os.Signal, 1)
	signal.Notify(quitCh, os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	<-quitCh

	// stop an application
	if err := a.Stop(context.Background()); err != nil {
		panic(err)
	}
}
