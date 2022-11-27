package broker

import (
	"context"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	banktypes "github.com/cosmos/cosmos-sdk/types"
)

func (b *Broker) PublishBank(ctx context.Context, response banktypes.Coins) error {

	jsonBytes, err := response.MarshalJSON()
	if err != nil {
		b.log.Error().Err(err).Msg("json marshal fail")
		return err
	}

	//err = b.writer.WriteMessages(ctx,
	//	// NOTE: Each Message has Topic defined, otherwise an error is returned.
	//	kafka.Message{
	//		Topic: "block-module",
	//		Key:   []byte("block"),
	//		Value: jsonBytes,
	//	},
	//)
	//if err != nil {
	//	b.log.Error().Err(err).Msg("failed to write messages")
	//	return err
	//}
	t := "test-topic"
	err = b.p.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &t, Partition: kafka.PartitionAny},
		Value:          jsonBytes,
		Headers:        []kafka.Header{{Key: "myTestHeader", Value: []byte("header values are binary")}},
	}, nil)
	if err != nil {
		return err
	}

	return nil
}
