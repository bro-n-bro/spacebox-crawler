package broker

import (
	"context"

	"github.com/pkg/errors"

	"bro-n-bro-osmosis/adapter/broker/model"

	jsoniter "github.com/json-iterator/go"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

var (
	UnbondingDelegationMessageTopic = "unbonding_delegation_message"
)

func (b *Broker) PublishUnbondingDelegationMessage(ctx context.Context, udm model.UnbondingDelegationMessage) error {
	return nil

	data, err := jsoniter.Marshal(udm) // FIXME: maybe user another way to encode data
	if err != nil {
		return errors.Wrap(err, MsgErrJsonMarshalFail)
	}
	err = b.p.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &UnbondingDelegationMessageTopic, Partition: kafka.PartitionAny},
		Value:          data,
		//Headers:        []kafka.Header{{Key: "myTestHeader", Value: []byte("header values are binary")}},
	}, nil)
	if err != nil {
		return errors.Wrap(err, "produce unbonding_delegation_message fail")
	}
	return nil
}
