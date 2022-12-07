package broker

import (
	"context"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"

	"bro-n-bro-osmosis/adapter/broker/model"
)

func (b *Broker) PublishBlock(ctx context.Context, block model.Block) error {
	return nil

	data, err := jsoniter.Marshal(block) // FIXME: maybe user another way to encode data
	if err != nil {
		return errors.Wrap(err, MsgErrJsonMarshalFail)
	}
	err = b.p.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: BlockTopic, Partition: kafka.PartitionAny},
		Value:          data,
		//Headers:        []kafka.Header{{Key: "myTestHeader", Value: []byte("header values are binary")}},
	}, nil)
	if err != nil {
		return errors.Wrap(err, "produce block fail")
	}

	return nil
}
