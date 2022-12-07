package broker

import (
	"context"

	"github.com/pkg/errors"

	"bro-n-bro-osmosis/adapter/broker/model"

	jsoniter "github.com/json-iterator/go"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

var ()

func (b *Broker) PublishMessage(ctx context.Context, message model.Message) error {
	return nil

	data, err := jsoniter.Marshal(message) // FIXME: maybe user another way to encode data
	if err != nil {
		return errors.Wrap(err, MsgErrJsonMarshalFail)
	}
	err = b.p.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: MessageTopic, Partition: kafka.PartitionAny},
		Value:          data,
		//Headers:        []kafka.Header{{Key: "myTestHeader", Value: []byte("header values are binary")}},
	}, nil)
	if err != nil {
		return errors.Wrap(err, "produce message fail")
	}
	return nil
}
