package storage

import (
	"context"

	"github.com/rs/zerolog"
)

type Storage struct {
	log *zerolog.Logger
	cfg Config
}

func (s *Storage) Start(ctx context.Context) error {
	return nil
}

func (s *Storage) Stop(ctx context.Context) error {
	return nil
}
