package broker

import (
	"context"

	"github.com/hexy-dev/spacebox/broker/model"
	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"
)

func (b *Broker) PublishUnbondingDelegation(ctx context.Context, ud model.UnbondingDelegation) error {

	data, err := jsoniter.Marshal(ud) // FIXME: maybe user another way to encode data
	if err != nil {
		return errors.Wrap(err, MsgErrJSONMarshalFail)
	}

	if err := b.produce(UnbondingDelegation, data); err != nil {
		return errors.Wrap(err, "produce unbonding_delegation_message fail")
	}
	return nil
}

func (b *Broker) PublishUnbondingDelegationMessage(ctx context.Context, udm model.UnbondingDelegationMessage) error {

	data, err := jsoniter.Marshal(udm) // FIXME: maybe user another way to encode data
	if err != nil {
		return errors.Wrap(err, MsgErrJSONMarshalFail)
	}

	if err := b.produce(UnbondingDelegationMessage, data); err != nil {
		return errors.Wrap(err, "produce unbonding_delegation_message fail")
	}
	return nil
}

func (b *Broker) PublishStakingParams(ctx context.Context, sp model.StakingParams) error {

	data, err := jsoniter.Marshal(sp) // FIXME: maybe user another way to encode data
	if err != nil {
		return errors.Wrap(err, MsgErrJSONMarshalFail)
	}

	if err := b.produce(StakingParams, data); err != nil {
		return errors.Wrap(err, "produce supply fail")
	}
	return nil
}

func (b *Broker) PublishDelegation(ctx context.Context, d model.Delegation) error {

	data, err := jsoniter.Marshal(d) // FIXME: maybe user another way to encode data
	if err != nil {
		return errors.Wrap(err, MsgErrJSONMarshalFail)
	}

	if err := b.produce(Delegation, data); err != nil {
		return err
	}
	return nil
}

func (b *Broker) PublishDelegationMessage(ctx context.Context, dm model.DelegationMessage) error {

	data, err := jsoniter.Marshal(dm) // FIXME: maybe user another way to encode data
	if err != nil {
		return errors.Wrap(err, MsgErrJSONMarshalFail)
	}

	if err := b.produce(DelegationMessage, data); err != nil {
		return errors.Wrap(err, "produce supply fail")
	}
	return nil
}
