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

	return &Broker{
		//writer: w,
		log: &l,
		cfg: cfg,
	}
}

func (b *Broker) Start(ctx context.Context) error {
	if !b.cfg.Enabled {
		return nil
	}

	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": b.cfg.ServerURL})

	if err != nil {
		fmt.Printf("Failed to create producer: %s\n", err)
		return err
	}
	b.p = p

	return nil
}

func (b *Broker) Stop(ctx context.Context) error {
	if !b.cfg.Enabled {
		return nil
	}
	b.p.Close()
	//return b.writer.Close()
	return nil
}

func (b *Broker) produce(topic Topic, data []byte) error {
	if !b.cfg.Enabled {
		return nil
	}
	err := b.p.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: topic, Partition: kafka.PartitionAny},
		Value:          data,
		//Headers:        []kafka.Header{{Key: "myTestHeader", Value: []byte("header values are binary")}},
	}, nil)
	if err != nil {
		return err
	}
	return nil
}
