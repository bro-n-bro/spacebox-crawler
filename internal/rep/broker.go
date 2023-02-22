package rep

import (
	"context"

	"github.com/bro-n-bro/spacebox/broker/model"
)

type Broker interface {
	// auth
	PublishAccount(context.Context, model.Account) error

	// core
	PublishBlock(context.Context, model.Block) error
	PublishMessage(ctx context.Context, message model.Message) error
	PublishTransaction(ctx context.Context, tx model.Transaction) error

	// bank
	PublishSupply(context.Context, model.Supply) error
	PublishSendMessage(context.Context, model.SendMessage) error
	PublishMultiSendMessage(ctx context.Context, msm model.MultiSendMessage) error
	PublishAccountBalance(ctx context.Context, ab model.AccountBalance) error

	// distribution
	PublishDelegationReward(context.Context, model.DelegationReward) error
	PublishDelegationRewardMessage(context.Context, model.DelegationRewardMessage) error
	PublishDistributionParams(ctx context.Context, dp model.DistributionParams) error
	PublishValidatorCommission(ctx context.Context, commission model.ValidatorCommission) error

	// staking
	PublishCommunityPool(ctx context.Context, cp model.CommunityPool) error
	PublishUnbondingDelegation(context.Context, model.UnbondingDelegation) error
	PublishUnbondingDelegationMessage(context.Context, model.UnbondingDelegationMessage) error
	PublishStakingParams(ctx context.Context, sp model.StakingParams) error
	PublishDelegation(ctx context.Context, d model.Delegation) error
	PublishDelegationMessage(ctx context.Context, dm model.DelegationMessage) error
	PublishRedelegationMessage(context.Context, model.RedelegationMessage) error
	PublishRedelegation(context.Context, model.Redelegation) error
	PublishStakingPool(ctx context.Context, sp model.StakingPool) error
	PublishValidator(ctx context.Context, val model.Validator) error
	PublishValidatorInfo(ctx context.Context, info model.ValidatorInfo) error
	PublishValidatorStatus(ctx context.Context, status model.ValidatorStatus) error
	PublishValidatorDescription(ctx context.Context, description model.ValidatorDescription) error

	// mint module
	PublishMintParams(ctx context.Context, mp model.MintParams) error
	PublishAnnualProvision(ctx context.Context, ap model.AnnualProvision) error

	// gov module
	PublishProposal(ctx context.Context, proposal model.Proposal) error
	PublishGovParams(ctx context.Context, params model.GovParams) error
	PublishProposalDeposit(ctx context.Context, pvm model.ProposalDeposit) error
	PublishProposalDepositMessage(ctx context.Context, pvm model.ProposalDepositMessage) error
	PublishProposalVoteMessage(context.Context, model.ProposalVoteMessage) error
	PublishProposalTallyResult(ctx context.Context, ptr model.ProposalTallyResult) error
}
