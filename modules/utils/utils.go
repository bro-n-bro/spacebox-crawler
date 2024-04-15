package utils

import (
	"os"

	"github.com/rs/zerolog"
)

func NewModuleLogger(name string) *zerolog.Logger {
	logger := zerolog.
		New(os.Stderr).
		Output(zerolog.ConsoleWriter{Out: os.Stderr}).
		With().Timestamp().
		Str("module", name).
		Logger()

	return &logger
}
