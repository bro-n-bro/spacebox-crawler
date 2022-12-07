package broker

import (
	"context"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"

	"bro-n-bro-osmosis/adapter/broker/model"
)

func (b *Broker) PublishDistributionParams(ctx context.Context, dp model.DistributionParams) error {
	return nil

	data, err := jsoniter.Marshal(dp) // FIXME: maybe user another way to encode data
	if err != nil {
		return errors.Wrap(err, MsgErrJsonMarshalFail)
	}
	err = b.p.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: DistributionParamsTopic, Partition: kafka.PartitionAny},
		Value:          data,
		//Headers:        []kafka.Header{{Key: "myTestHeader", Value: []byte("header values are binary")}},
	}, nil)
	if err != nil {
		return errors.Wrap(err, "produce block fail")
	}

	return nil
}
