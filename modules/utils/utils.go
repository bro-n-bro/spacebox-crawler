package utils

import (
	"encoding/base64"
	"os"

	"github.com/rs/zerolog"
	"golang.org/x/exp/constraints"
)

func ContainAny[T constraints.Ordered](src []T, trg T) bool {
	for _, v := range src {
		if v == trg {
			return true
		}
	}

	return false
}

func DecodeToString(v string) (string, error) {
	val, err := base64.StdEncoding.DecodeString(v)
	if err != nil {
		return "", err
	}

	return string(val), nil
}
func NewModuleLogger(name string) *zerolog.Logger {
	logger := zerolog.
		New(os.Stderr).
		Output(zerolog.ConsoleWriter{Out: os.Stderr}).
		With().Timestamp().
		Str("module", name).
		Logger()

	return &logger
}
