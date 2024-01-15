package broker

import (
	"context"
	"fmt"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

const (
	MsgErrJSONMarshalFail   = "json marshal fail: %w"
	MsgErrCreateProducer    = "can't create producer connection to broker: %w "
	MsgErrCreateAdminClient = "can't create admin client connection to broker: %w"
	MsgErrCreateTopics      = "can't create topics in broker: %w"
	MsgErrProduceTopic      = "can't produce topic: %w"
	MsgErrCreatePartitions  = "can't create partitions in broker: %w"
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
		validator      cache[string, int64]
		valCommission  cache[string, int64]
		valDescription cache[string, int64]
		valInfo        cache[string, int64]
		valStatus      cache[string, int64]
	}

	cache[K, V comparable] interface {
		UpdateCacheValue(K, V) bool
	}

	opt func(b *Broker)
)

func New(cfg Config, modules []string, l zerolog.Logger, opts ...opt) *Broker {
	l = l.With().Str("cmp", "broker").Logger()

	b := &Broker{
		log:     &l,
		cfg:     cfg,
		modules: modules,
	}

	for _, apply := range opts {
		apply(b)
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
		// "message.max.bytes": 5 << 20, // 5 MB
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
				b.log.Error().Str("topic_partition", m.TopicPartition.String()).Err(err).Msg("delivery error")
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

// marshalAndProduce marshals the message to JSON and produces it to the kafka.
func (b *Broker) marshalAndProduce(topic Topic, msg interface{}) error {
	data, err := jsoniter.Marshal(msg)
	if err != nil {
		return errors.Wrap(err, MsgErrJSONMarshalFail)
	}

	if err = b.produce(topic, data); err != nil {
		return errors.Wrap(err, MsgErrProduceTopic)
	}

	return nil
}

// produce produces the message to the kafka.
func (b *Broker) produce(topic Topic, data []byte) error {
	if !b.cfg.Enabled {
		return nil
	}

	err := b.p.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: topic, Partition: kafka.PartitionAny},
		Value:          data,
	}, nil)

	if kafkaError, ok := err.(kafka.Error); ok && kafkaError.Code() == kafka.ErrQueueFull {
		b.log.Info().Str("topic", *topic).Msg("kafka local queue full error. Going to Flush then retry")
		flushedMessages := b.p.Flush(30 * 1000)
		b.log.Info().Str("topic", *topic).Int("flushed_messages", flushedMessages).
			Msg("flushed kafka messages. Outstanding events still un-flushed")

		return b.produce(topic, data)
	}

	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("produce %s fail", *topic))
	}

	return nil
}

// getCurrentTopics returns the list of topics based on enabled modules.
// nolint:gocyclo
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
		case "liquidity":
			topics = append(topics, liquidityTopics.ToStringSlice()...)
		case "graph":
			topics = append(topics, graphTopics.ToStringSlice()...)
		case "bandwidth":
			topics = append(topics, bandwidthTopics.ToStringSlice()...)
		case "dmn":
			topics = append(topics, dmnTopics.ToStringSlice()...)
		case "grid":
			topics = append(topics, gridTopics.ToStringSlice()...)
		case "rank":
			topics = append(topics, rankTopics.ToStringSlice()...)
		case "resources":
			topics = append(topics, resourcesTopics.ToStringSlice()...)
		case "wasm":
			topics = append(topics, wasmTopics.ToStringSlice()...)
		case "raw":
			topics = append(topics, rawTopics.ToStringSlice()...)
		default:
			b.log.Warn().Str("name", m).Msg("unknown module in config")
			continue
		}
	}

	return removeDuplicates(topics)
}

func WithValidatorCache(valCache cache[string, int64]) func(b *Broker) {
	return func(b *Broker) {
		b.cache.validator = valCache
	}
}

func WithValidatorCommissionCache(valCommissionCache cache[string, int64]) func(b *Broker) {
	return func(b *Broker) {
		b.cache.valCommission = valCommissionCache
	}
}

func WithValidatorDescriptionCache(valDescriptionCache cache[string, int64]) func(b *Broker) {
	return func(b *Broker) {
		b.cache.valDescription = valDescriptionCache
	}
}

func WithValidatorInfoCache(valInfoCache cache[string, int64]) func(b *Broker) {
	return func(b *Broker) {
		b.cache.valInfo = valInfoCache
	}
}

func WithValidatorStatusCache(valStatusCache cache[string, int64]) func(b *Broker) {
	return func(b *Broker) {
		b.cache.valStatus = valStatusCache
	}
}

func removeDuplicates[T comparable](s []T) []T {
	res := make([]T, 0)
	uniq := make(map[T]struct{})

	for i := 0; i < len(s); i++ {
		if _, ok := uniq[s[i]]; !ok {
			uniq[s[i]] = struct{}{}
			res = append(res, s[i])
		}
	}

	return res
}
