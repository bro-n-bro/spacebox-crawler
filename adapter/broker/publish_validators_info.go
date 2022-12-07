package broker

import (
	"context"

	"github.com/pkg/errors"

	"bro-n-bro-osmosis/adapter/broker/model"

	jsoniter "github.com/json-iterator/go"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func (b *Broker) PublishValidatorsInfo(ctx context.Context, infos []model.ValidatorInfo) error {
	return nil

	for i := 0; i < len(infos); i++ {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
		data, err := jsoniter.Marshal(infos[i]) // FIXME: maybe user another way to encode data
		if err != nil {
			return errors.Wrap(err, MsgErrJsonMarshalFail)
		}
		err = b.p.Produce(&kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: ValidatorInfo, Partition: kafka.PartitionAny},
			Value:          data,
			//Headers:        []kafka.Header{{Key: "myTestHeader", Value: []byte("header values are binary")}},
		}, nil)
		if err != nil {
			return errors.Wrap(err, "produce account fail")
		}
	}
	return nil
}
