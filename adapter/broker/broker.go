package broker

import (
	"context"
	"os"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

const (
	MsgErrJsonMarshalFail   = "json marshal fail: %w"
	MsgErrCreateProducer    = "cant create producer connection to broker: %w "
	MsgErrCreateAdminClient = "cant create admin client connection to broker: %w"
	MsgErrCreateTopics      = "cant create topics in broker: %w"
)

type Broker struct {
	log     *zerolog.Logger
	p       *kafka.Producer
	ac      *kafka.AdminClient
	cfg     Config
	modules []string
}

func New(cfg Config, modules []string) *Broker {
	l := zerolog.New(os.Stderr).Output(zerolog.ConsoleWriter{Out: os.Stderr}).With().Timestamp().
		Str("cmp", "broker").Logger()

	return &Broker{
		log:     &l,
		cfg:     cfg,
		modules: modules,
	}
}

func (b *Broker) Start(ctx context.Context) error {
	if !b.cfg.Enabled {
		return nil
	}

	// create an admin client connection
	ac, err := kafka.NewAdminClient(&kafka.ConfigMap{"bootstrap.servers": b.cfg.ServerURL})
	if err != nil {
		b.log.Error().Err(err).Msgf(MsgErrCreateAdminClient, err)
		return errors.Wrap(err, MsgErrCreateAdminClient)
	}

	// get enabled topics based on enabled modules
	topics := b.getCurrentTopics(b.modules)
	kafkaTopics := make([]kafka.TopicSpecification, len(topics))
	for i, topic := range topics {
		kafkaTopics[i] = kafka.TopicSpecification{
			Topic:         topic,
			NumPartitions: 1,
		}
	}

	// create init topics if needed
	_, err = ac.CreateTopics(ctx, kafkaTopics)
	if err != nil {
		b.log.Error().Err(err).Msgf(MsgErrCreateTopics, err)
		return errors.Wrap(err, MsgErrCreateTopics)
	}

	// create a producer connection
	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": b.cfg.ServerURL})
	if err != nil {
		b.log.Error().Err(err).Msgf(MsgErrCreateProducer, err)
		return errors.New(MsgErrCreateProducer)
	}

	b.p = p
	b.ac = ac

	b.log.Info().Msg("broker started")
	return nil
}

func (b *Broker) Stop(ctx context.Context) error {
	if !b.cfg.Enabled {
		return nil
	}
	b.p.Close()
	b.ac.Close()
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

func (b *Broker) getCurrentTopics(modules []string) []string {
	topics := make([]string, 0)
	for _, m := range modules {
		switch m {
		case "auth":
			topics = append(topics, authTopics.ToStringSlice()...)
		case "bank":
			topics = append(topics, bankTopics.ToStringSlice()...)
		case "gov":
			topics = append(topics, govTopics.ToStringSlice()...)
		case "mint":
			topics = append(topics, mintTopics.ToStringSlice()...)
		case "staking":
			topics = append(topics, stakingTopics.ToStringSlice()...)
		case "distribution":
			topics = append(topics, distributionTopics.ToStringSlice()...)
		case "core":
			topics = append(topics, coreTopics.ToStringSlice()...)
		default:
			b.log.Warn().Msgf("unknown module in config: %v", m)
			continue
		}
	}
	return topics
}
