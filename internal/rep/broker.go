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
	PublishValidatorPrecommit(ctx context.Context, vp model.ValidatorPrecommit) error
	PublishValidatorVotingPower(ctx context.Context, vp model.ValidatorVotingPower) error

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
	PublishSetWithdrawAddressMessage(_ context.Context, swm model.SetWithdrawAddressMessage) error
	PublishProposerReward(ctx context.Context, pr model.ProposerReward) error
	PublishDistributionCommission(ctx context.Context, commission model.DistributionCommission) error
	PublishDistributionReward(ctx context.Context, reward model.DistributionReward) error
	PublishWithdrawValidatorCommissionMessage(_ context.Context, wvcm model.WithdrawValidatorCommissionMessage) error

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
	PublishCreateValidatorMessage(ctx context.Context, cvm model.CreateValidatorMessage) error
	PublishEditValidatorMessage(ctx context.Context, message model.EditValidatorMessage) error
	PublishCancelUnbondingDelegationMessage(_ context.Context, description model.CancelUnbondingDelegationMessage) error

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
	PublishSubmitProposalMessage(ctx context.Context, spm model.SubmitProposalMessage) error
	PublishVoteWeightedMessage(ctx context.Context, vwm model.VoteWeightedMessage) error

	// authz module
	PublishGrantMessage(context.Context, model.GrantMessage) error
	PublishAuthzGrant(context.Context, model.AuthzGrant) error
	PublishRevokeMessage(context.Context, model.RevokeMessage) error
	PublishExecMessage(context.Context, model.ExecMessage) error

	// feegrant module
	PublishFeeAllowance(context.Context, model.FeeAllowance) error
	PublishGrantAllowanceMessage(context.Context, model.GrantAllowanceMessage) error
	PublishRevokeAllowanceMessage(context.Context, model.RevokeAllowanceMessage) error

	// slashing module
	PublishSlashingParams(context.Context, model.SlashingParams) error
	PublishUnjailMessage(context.Context, model.UnjailMessage) error
	PublishHandleValidatorSignature(ctx context.Context, msg model.HandleValidatorSignature) error

	// ibc
	PublishTransferMessage(context.Context, model.TransferMessage) error
	PublishAcknowledgementMessage(context.Context, model.AcknowledgementMessage) error
	PublishReceivePacketMessage(context.Context, model.RecvPacketMessage) error
	PublishDenomTrace(context.Context, model.DenomTrace) error

	// liquidity
	PublishSwap(context.Context, model.Swap) error
	PublishLiquidityPool(context.Context, model.LiquidityPool) error
}
