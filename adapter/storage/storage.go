package storage

import (
	"context"

	mongoprom "github.com/globocom/mongo-go-prometheus"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Storage struct {
	log        *zerolog.Logger
	cli        *mongo.Client
	collection *mongo.Collection

	cfg Config
}

func New(cfg Config, l zerolog.Logger) *Storage {
	l = l.With().Str("cmp", "mongo").Logger()

	return &Storage{
		cfg: cfg,
		log: &l,
	}
}

func (s *Storage) Start(ctx context.Context) error {
	opts := []*options.ClientOptions{
		options.Client().ApplyURI(s.cfg.URI),
		options.Client().SetMaxPoolSize(s.cfg.MaxPoolSize),
		options.Client().SetMaxConnecting(s.cfg.MaxConnecting),
		options.Client().SetAuth(options.Credential{
			Username: s.cfg.User,
			Password: s.cfg.Password,
		}),
	}

	if s.cfg.MetricsEnabled {
		monitor := mongoprom.NewCommandMonitor(
			mongoprom.WithInstanceName("blocks"),
			mongoprom.WithDurationBuckets([]float64{.001, .005, .01}),
		)
		opts = append(opts, options.Client().SetMonitor(monitor))
	}

	client, err := mongo.Connect(ctx, opts...)
	if err != nil {
		return err
	}

	s.cli = client

	if err := s.Ping(ctx); err != nil {
		return err
	}

	collection := s.cli.Database("spacebox").Collection("blocks")
	s.collection = collection

	mod := mongo.IndexModel{
		Keys: bson.M{"height": 1}, // index in ascending order or -1 for descending order
		// Options: options.Index().SetUnique(true),
	}

	if _, err := collection.Indexes().CreateOne(ctx, mod); err != nil {
		return err
	}

	return nil
}

func (s *Storage) Stop(ctx context.Context) error {
	s.log.Info().Msg("start setErrorStatusForProcessing")

	// TODO:
	if err := s.setErrorStatusForProcessing(ctx); err != nil {
		s.log.Error().Err(err).Msg("setErrorStatusForProcessing error")
		return err
	}

	return s.cli.Disconnect(ctx)
}

func (s *Storage) Ping(ctx context.Context) error {
	return s.cli.Ping(ctx, nil)
}
