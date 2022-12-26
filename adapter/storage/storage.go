package storage

import (
	"context"
	"os"

	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Storage struct {
	log        *zerolog.Logger
	cli        *mongo.Client
	collection *mongo.Collection
	cfg        Config
}

func New(cfg Config) *Storage {
	l := zerolog.New(os.Stderr).Output(zerolog.ConsoleWriter{Out: os.Stderr}).With().Timestamp().
		Str("cmp", "mongo").Logger()

	return &Storage{
		cfg: cfg,
		log: &l,
	}
}

func (s *Storage) Start(ctx context.Context) error {
	opts := []*options.ClientOptions{
		options.Client().ApplyURI(s.cfg.URI),
		options.Client().SetMaxPoolSize(8),
		options.Client().SetMaxConnecting(8),
		options.Client().SetAuth(options.Credential{
			AuthMechanism:           "",
			AuthMechanismProperties: nil,
			AuthSource:              "",
			Username:                s.cfg.User,
			Password:                s.cfg.Password,
			PasswordSet:             false,
		}),
	}

	client, err := mongo.Connect(ctx, opts...)
	if err != nil {
		return err
	}
	s.cli = client

	if err := s.cli.Ping(ctx, nil); err != nil {
		return err
	}

	collection := s.cli.Database("spacebox").Collection("blocks")
	s.collection = collection

	s.log.Info().Msg("storage started")

	return nil
}

func (s *Storage) Stop(ctx context.Context) error {
	s.log.Info().Msg("start setErrorStatusForProcessing")

	// TODO:
	err := s.setErrorStatusForProcessing(context.Background())
	if err != nil {
		s.log.Error().Err(err).Msg("setErrorStatusForProcessing error")
		return err
	}
	return s.cli.Disconnect(ctx)
}
