package broker

import (
	"bro-n-bro-osmosis/adapter/broker/model"
	"context"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"
)

func (b *Broker) PublishUnbondingDelegation(ctx context.Context, ud model.UnbondingDelegation) error {
	return nil

	data, err := jsoniter.Marshal(ud) // FIXME: maybe user another way to encode data
	if err != nil {
		return errors.Wrap(err, MsgErrJsonMarshalFail)
	}
	err = b.p.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: UnbondingDelegation, Partition: kafka.PartitionAny},
		Value:          data,
		//Headers:        []kafka.Header{{Key: "myTestHeader", Value: []byte("header values are binary")}},
	}, nil)
	if err != nil {
		return errors.Wrap(err, "produce unbonding_delegation_message fail")
	}
	return nil
}

func (b *Broker) PublishUnbondingDelegationMessage(ctx context.Context, udm model.UnbondingDelegationMessage) error {
	return nil

	data, err := jsoniter.Marshal(udm) // FIXME: maybe user another way to encode data
	if err != nil {
		return errors.Wrap(err, MsgErrJsonMarshalFail)
	}
	err = b.p.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: UnbondingDelegationMessage, Partition: kafka.PartitionAny},
		Value:          data,
		//Headers:        []kafka.Header{{Key: "myTestHeader", Value: []byte("header values are binary")}},
	}, nil)
	if err != nil {
		return errors.Wrap(err, "produce unbonding_delegation_message fail")
	}
	return nil
}

func (b *Broker) PublishStakingParams(ctx context.Context, sp model.StakingParams) error {
	return nil

	data, err := jsoniter.Marshal(sp) // FIXME: maybe user another way to encode data
	if err != nil {
		return errors.Wrap(err, MsgErrJsonMarshalFail)
	}
	err = b.p.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: StakingParams, Partition: kafka.PartitionAny},
		Value:          data,
		//Headers:        []kafka.Header{{Key: "myTestHeader", Value: []byte("header values are binary")}},
	}, nil)
	if err != nil {
		return errors.Wrap(err, "produce supply fail")
	}
	return nil
}

func (b *Broker) PublishDelegation(ctx context.Context, d model.Delegation) error {
	return nil

	data, err := jsoniter.Marshal(d) // FIXME: maybe user another way to encode data
	if err != nil {
		return errors.Wrap(err, MsgErrJsonMarshalFail)
	}
	err = b.p.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: Delegation, Partition: kafka.PartitionAny},
		Value:          data,
		//Headers:        []kafka.Header{{Key: "myTestHeader", Value: []byte("header values are binary")}},
	}, nil)
	if err != nil {
		return errors.Wrap(err, "produce supply fail")
	}
	return nil
}

func (b *Broker) PublishDelegationMessage(ctx context.Context, dm model.DelegationMessage) error {
	return nil

	data, err := jsoniter.Marshal(dm) // FIXME: maybe user another way to encode data
	if err != nil {
		return errors.Wrap(err, MsgErrJsonMarshalFail)
	}
	err = b.p.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: DelegationMessage, Partition: kafka.PartitionAny},
		Value:          data,
		//Headers:        []kafka.Header{{Key: "myTestHeader", Value: []byte("header values are binary")}},
	}, nil)
	if err != nil {
		return errors.Wrap(err, "produce supply fail")
	}
	return nil
}
