package main

import (
	"log"
	"os"
	"time"

	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"

	"github.com/bro-n-bro/spacebox-crawler/internal/app"
	executor "github.com/bro-n-bro/spacebox-crawler/pkg/app"
)

const (
	DefaultEnvFile = ".env"
	EnvFile        = "ENV_FILE"
)

// Version provided by ldflags
var Version = "develop"

func main() {
	// try to get .env file from Environments
	fileName, ok := os.LookupEnv(EnvFile)
	if !ok {
		fileName = DefaultEnvFile
	}

	// load environment variables based on .env file
	if err := godotenv.Load(fileName); err != nil {
		log.Fatal(err)
	}

	var cfg app.Config
	// fill these variables into a struct
	if err := env.Parse(&cfg); err != nil {
		log.Fatal(err)
	}

	// parse log level
	logLevel, err := zerolog.ParseLevel(cfg.LogLevel)
	if err != nil {
		log.Fatal(err)
	}

	// create a logger instance
	// TODO: log in file
	logger := zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.DateTime}).
		Level(logLevel).
		Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.DateTime}).
		With().
		Timestamp().
		Logger()

	// create an application
	a := app.New(cfg, Version, logger)

	// run service
	if err := executor.Run(a); err != nil {
		log.Fatal(err)
	}
}
