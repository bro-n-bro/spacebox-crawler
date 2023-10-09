package staking

import (
	"context"

	"github.com/bro-n-bro/spacebox/broker/model"
)

type broker interface {
	PublishAccount(context.Context, model.Account) error // FIXME: method from auth module

	PublishUnbondingDelegation(context.Context, model.UnbondingDelegation) error
	PublishUnbondingDelegationMessage(context.Context, model.UnbondingDelegationMessage) error
	PublishStakingParams(ctx context.Context, sp model.StakingParams) error
	PublishDelegation(ctx context.Context, d model.Delegation) error
	PublishDisabledDelegation(ctx context.Context, d model.Delegation) error
	PublishDelegationMessage(ctx context.Context, dm model.DelegationMessage) error
	PublishStakingPool(ctx context.Context, sp model.StakingPool) error
	PublishValidator(ctx context.Context, val model.Validator) error
	PublishValidatorInfo(ctx context.Context, infos model.ValidatorInfo) error
	PublishValidatorStatus(ctx context.Context, statuses model.ValidatorStatus) error
	PublishValidatorDescription(ctx context.Context, description model.ValidatorDescription) error
	PublishValidatorCommission(ctx context.Context, commission model.ValidatorCommission) error
	PublishRedelegation(context.Context, model.Redelegation) error
	PublishRedelegationMessage(context.Context, model.RedelegationMessage) error
	PublishCreateValidatorMessage(ctx context.Context, cvm model.CreateValidatorMessage) error
	PublishEditValidatorMessage(ctx context.Context, message model.EditValidatorMessage) error
	PublishCancelUnbondingDelegationMessage(_ context.Context, description model.CancelUnbondingDelegationMessage) error
}
