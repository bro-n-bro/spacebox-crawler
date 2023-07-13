package broker

import (
	"context"
	"fmt"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

const (
	MsgErrJSONMarshalFail   = "json marshal fail: %w"
	MsgErrCreateProducer    = "cant create producer connection to broker: %w "
	MsgErrCreateAdminClient = "cant create admin client connection to broker: %w"
	MsgErrCreateTopics      = "cant create topics in broker: %w"
	MsgErrCreatePartitions  = "cant create partitions in broker: %w"
)

type (
	Broker struct {
		log     *zerolog.Logger
		p       *kafka.Producer
		ac      *kafka.AdminClient
		cache   lruCache
		modules []string
		cfg     Config
	}

	lruCache struct {
		account        cacheI[string, int64]
		validator      cacheI[string, int64]
		valCommission  cacheI[string, int64]
		valDescription cacheI[string, int64]
		valInfo        cacheI[string, int64]
		valStatus      cacheI[string, int64]
	}

	cacheI[K, V comparable] interface {
		UpdateCacheValue(K, V) bool
	}

	opts func(b *Broker)
)

func New(cfg Config, modules []string, l zerolog.Logger, opts ...opts) *Broker {
	l = l.With().Str("cmp", "broker").Logger()

	b := &Broker{
		log:     &l,
		cfg:     cfg,
		modules: modules,
	}

	for _, opt := range opts {
		opt(b)
	}

	return b
}

func (b *Broker) Start(ctx context.Context) error {
	if !b.cfg.Enabled {
		return nil
	}

	// create an admin client connection
	ac, err := kafka.NewAdminClient(&kafka.ConfigMap{
		"bootstrap.servers": b.cfg.ServerURL,
	})
	if err != nil {
		b.log.Error().Err(err).Msg(MsgErrCreateAdminClient)
		return errors.Wrap(err, MsgErrCreateAdminClient)
	}

	// get enabled topics based on enabled modules
	topics := b.getCurrentTopics(b.modules)
	kafkaTopics := make([]kafka.TopicSpecification, len(topics))
	// kafkaPartitions := make([]kafka.PartitionsSpecification, len(topics))
	for i, topic := range topics {
		kafkaTopics[i] = kafka.TopicSpecification{
			Topic:         topic,
			NumPartitions: b.cfg.PartitionsCount,
		}
		// kafkaPartitions[i] = kafka.PartitionsSpecification{
		//	Topic:      topic,
		//	IncreaseTo: 2,
		// }
	}

	// create init topics if needed
	if _, err = ac.CreateTopics(ctx, kafkaTopics); err != nil {
		b.log.Error().Err(err).Msg(MsgErrCreateTopics)
		return errors.Wrap(err, MsgErrCreateTopics)
	}

	// create a producer connection
	p, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": b.cfg.ServerURL,
		"message.max.bytes": 5 << 20, // 5 MB
	})
	if err != nil {
		b.log.Error().Err(err).Msg(MsgErrCreateProducer)
		return errors.New(MsgErrCreateProducer)
	}

	go func(drs chan kafka.Event) {
		for ev := range drs {
			m, ok := ev.(*kafka.Message)
			if !ok {
				continue
			}

			if err := m.TopicPartition.Error; err != nil {
				b.log.Error().Err(err).Msgf("Delivery error: %v", m.TopicPartition)
			}
		}
	}(p.Events())

	b.p = p
	b.ac = ac

	return nil
}

func (b *Broker) Stop(ctx context.Context) error {
	if !b.cfg.Enabled {
		return nil
	}

	b.p.Flush(30 * 1000)

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
	}, nil)

	if kafkaError, ok := err.(kafka.Error); ok && kafkaError.Code() == kafka.ErrQueueFull {
		b.log.Info().Str("topic", *topic).Msg("Kafka local queue full error - Going to Flush then retry...")
		flushedMessages := b.p.Flush(30 * 1000)
		b.log.Info().Str("topic", *topic).
			Msgf("Flushed kafka messages. Outstanding events still un-flushed: %d", flushedMessages)

		return b.produce(topic, data)
	}

	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("produce %s fail", *topic))
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
		case "authz":
			topics = append(topics, authzTopics.ToStringSlice()...)
		case "feegrant":
			topics = append(topics, feegrantTopics.ToStringSlice()...)
		case "slashing":
			topics = append(topics, slashingTopics.ToStringSlice()...)
		case "ibc":
			topics = append(topics, ibcTopics.ToStringSlice()...)
		default:
			b.log.Warn().Msgf("unknown module in config: %v", m)
			continue
		}
	}

	return topics
}

func WithAccountCache(accCache cacheI[string, int64]) func(b *Broker) {
	return func(b *Broker) {
		b.cache.account = accCache
	}
}

func WithValidatorCache(valCache cacheI[string, int64]) func(b *Broker) {
	return func(b *Broker) {
		b.cache.validator = valCache
	}
}

func WithValidatorCommissionCache(valCommissionCache cacheI[string, int64]) func(b *Broker) {
	return func(b *Broker) {
		b.cache.valCommission = valCommissionCache
	}
}

func WithValidatorDescriptionCache(valDescriptionCache cacheI[string, int64]) func(b *Broker) {
	return func(b *Broker) {
		b.cache.valDescription = valDescriptionCache
	}
}

func WithValidatorInfoCache(valInfoCache cacheI[string, int64]) func(b *Broker) {
	return func(b *Broker) {
		b.cache.valInfo = valInfoCache
	}
}

func WithValidatorStatusCache(valStatusCache cacheI[string, int64]) func(b *Broker) {
	return func(b *Broker) {
		b.cache.valStatus = valStatusCache
	}
}
