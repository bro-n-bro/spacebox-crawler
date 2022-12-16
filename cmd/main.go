package main

import (
	"log"
	"os"

	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"

	"bro-n-bro-osmosis/internal/app"
	executor "bro-n-bro-osmosis/pkg/app"
)

const (
	DefaultEnvFile = ".env"
	EnvFile        = "ENV_FILE"
)

func main() {
	_main()
}

func _main() {
	// try to get .env file from Environments
	fileName, ok := os.LookupEnv(EnvFile)
	if !ok {
		fileName = DefaultEnvFile
	}

	// load environment variables based on .env file
	if err := godotenv.Load(fileName); err != nil {
		panic(err)
	}

	var cfg app.Config
	// fill these variables into a struct
	if err := env.Parse(&cfg); err != nil {
		panic(err)
	}

	// create an application
	a := app.New(cfg)

	// run service
	if err := executor.Run(a); err != nil {
		log.Fatal(err)
	}
}
