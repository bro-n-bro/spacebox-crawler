package rep

import (
	"context"

	"github.com/hexy-dev/spacebox/broker/model"
)

type Broker interface {
	PublishAccounts(context.Context, []model.Account) error
	PublishValidators(ctx context.Context, vals []model.Validator) error
	PublishValidatorsInfo(ctx context.Context, infos []model.ValidatorInfo) error
	PublishValidatorsStatuses(ctx context.Context, statuses []model.ValidatorStatus) error

	PublishBlock(context.Context, model.Block) error
	PublishCommunityPool(ctx context.Context, cp model.CommunityPool) error
	PublishSupply(context.Context, model.Supply) error
	PublishSendMessage(context.Context, model.SendMessage) error
	PublishDelegationRewardMessage(context.Context, model.DelegationRewardMessage) error
	PublishProposalVoteMessage(context.Context, model.ProposalVoteMessage) error
	PublishProposalTallyResult(ctx context.Context, ptr model.ProposalTallyResult) error
	PublishRedelegationMessage(context.Context, model.RedelegationMessage) error
	PublishRedelegation(context.Context, model.Redelegation) error
	PublishMultiSendMessage(ctx context.Context, msm model.MultiSendMessage) error
	PublishDistributionParams(ctx context.Context, dp model.DistributionParams) error

	PublishAccountBalance(ctx context.Context, ab model.AccountBalance) error

	// staking
	PublishUnbondingDelegation(context.Context, model.UnbondingDelegation) error
	PublishUnbondingDelegationMessage(context.Context, model.UnbondingDelegationMessage) error
	PublishStakingParams(ctx context.Context, sp model.StakingParams) error
	PublishDelegation(ctx context.Context, d model.Delegation) error
	PublishDelegationMessage(ctx context.Context, dm model.DelegationMessage) error
	PublishStakingPool(ctx context.Context, sp model.StakingPool) error

	// mint module
	PublishMintParams(ctx context.Context, mp model.MintParams) error
	PublishAnnualProvision(ctx context.Context, ap model.AnnualProvision) error

	// gov module
	PublishGovParams(ctx context.Context, params model.GovParams) error

	// block_chain custom module
	PublishMessage(ctx context.Context, message model.Message) error
	PublishTransaction(ctx context.Context, tx model.Transaction) error
}
