package broker

import (
	"context"
	"fmt"
	"os"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/rs/zerolog"
)

const (
	MsgErrJsonMarshalFail = "json marshal fail: %w"
)

type Broker struct {
	log *zerolog.Logger
	//writer *kafka.Writer
	p   *kafka.Producer
	cfg Config
}

func New(cfg Config) *Broker {
	l := zerolog.New(os.Stderr).Output(zerolog.ConsoleWriter{Out: os.Stderr}).With().Timestamp().
		Str("cmp", "broker").Logger()

	//w := &kafka.Writer{
	//	Addr:                   kafka.TCP("localhost:9092"),
	//	AllowAutoTopicCreation: true,
	//}

	return &Broker{
		//writer: w,
		log: &l,
		cfg: cfg,
	}
}

func (b *Broker) Start(ctx context.Context) error {
	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": b.cfg.ServerURL})

	if err != nil {
		fmt.Printf("Failed to create producer: %s\n", err)
		return err
	}
	b.p = p

	return nil
}

func (b *Broker) Stop(ctx context.Context) error {
	b.p.Close()
	//return b.writer.Close()
	return nil
}
