package main

import (
	"log"

	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"

	"bro-n-bro-osmosis/internal/app"
	executor "bro-n-bro-osmosis/pkg/app"
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

	// run service
	if err := executor.Run(a); err != nil {
		log.Fatal(err)
	}
}
