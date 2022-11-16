package broker

import (
	"context"
	"os"

	"github.com/rs/zerolog"
	"github.com/segmentio/kafka-go"
)

type Broker struct {
	log    *zerolog.Logger
	writer *kafka.Writer
}

func New() *Broker {
	l := zerolog.New(os.Stderr).Output(zerolog.ConsoleWriter{Out: os.Stderr}).With().Timestamp().
		Str("cmp", "broker").Logger()

	w := &kafka.Writer{
		Addr:                   kafka.TCP("localhost:9092"),
		AllowAutoTopicCreation: true,
	}

	return &Broker{
		writer: w,
		log:    &l,
	}
}

func (b *Broker) Start(ctx context.Context) error {
	return nil
}

func (b *Broker) Stop(ctx context.Context) error {
	return b.writer.Close()
}
