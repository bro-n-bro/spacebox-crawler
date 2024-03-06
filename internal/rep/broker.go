package rep

import (
	"context"

	"github.com/bro-n-bro/spacebox/broker/model"
)

type Broker interface {
	// auth module
	PublishAccount(context.Context, model.Account) error

	// core module
	PublishBlock(context.Context, model.Block) error
	PublishMessage(context.Context, model.Message) error
	PublishTransaction(context.Context, model.Transaction) error
	PublishValidatorPreCommit(context.Context, model.ValidatorPreCommit) error
	PublishValidatorVotingPower(context.Context, model.ValidatorVotingPower) error

	// bank module
	PublishSupply(context.Context, model.Supply) error
	PublishSendMessage(context.Context, model.SendMessage) error
	PublishMultiSendMessage(context.Context, model.MultiSendMessage) error
	PublishAccountBalance(context.Context, model.AccountBalance) error

	// distribution module
	PublishDelegationReward(context.Context, model.DelegationReward) error
	PublishDelegationRewardMessage(context.Context, model.DelegationRewardMessage) error
	PublishDistributionParams(context.Context, model.DistributionParams) error
	PublishValidatorCommission(context.Context, model.ValidatorCommission) error
	PublishSetWithdrawAddressMessage(context.Context, model.SetWithdrawAddressMessage) error
	PublishProposerReward(context.Context, model.ProposerReward) error
	PublishDistributionCommission(context.Context, model.DistributionCommission) error
	PublishDistributionReward(context.Context, model.DistributionReward) error
	PublishWithdrawValidatorCommissionMessage(context.Context, model.WithdrawValidatorCommissionMessage) error

	// staking module
	PublishCommunityPool(context.Context, model.CommunityPool) error
	PublishUnbondingDelegation(context.Context, model.UnbondingDelegation) error
	PublishUnbondingDelegationMessage(context.Context, model.UnbondingDelegationMessage) error
	PublishStakingParams(context.Context, model.StakingParams) error
	PublishDelegation(context.Context, model.Delegation) error
	PublishDisabledDelegation(context.Context, model.Delegation) error
	PublishDelegationMessage(context.Context, model.DelegationMessage) error
	PublishRedelegationMessage(context.Context, model.RedelegationMessage) error
	PublishRedelegation(context.Context, model.Redelegation) error
	PublishStakingPool(context.Context, model.StakingPool) error
	PublishValidator(context.Context, model.Validator) error
	PublishValidatorInfo(context.Context, model.ValidatorInfo) error
	PublishValidatorStatus(context.Context, model.ValidatorStatus) error
	PublishValidatorDescription(context.Context, model.ValidatorDescription) error
	PublishCreateValidatorMessage(context.Context, model.CreateValidatorMessage) error
	PublishEditValidatorMessage(context.Context, model.EditValidatorMessage) error
	PublishCancelUnbondingDelegationMessage(context.Context, model.CancelUnbondingDelegationMessage) error

	// mint module
	PublishMintParams(context.Context, model.MintParams) error
	PublishAnnualProvision(context.Context, model.AnnualProvision) error

	// gov module
	PublishProposal(context.Context, model.Proposal) error
	PublishGovParams(context.Context, model.GovParams) error
	PublishProposalDeposit(context.Context, model.ProposalDeposit) error
	PublishProposalDepositMessage(context.Context, model.ProposalDepositMessage) error
	PublishProposalVoteMessage(context.Context, model.ProposalVoteMessage) error
	PublishProposalTallyResult(context.Context, model.ProposalTallyResult) error
	PublishSubmitProposalMessage(context.Context, model.SubmitProposalMessage) error
	PublishVoteWeightedMessage(context.Context, model.VoteWeightedMessage) error

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

	// ibc module
	PublishTransferMessage(context.Context, model.TransferMessage) error
	PublishAcknowledgementMessage(context.Context, model.AcknowledgementMessage) error
	PublishReceivePacketMessage(context.Context, model.RecvPacketMessage) error
	PublishDenomTrace(context.Context, model.DenomTrace) error

	// liquidity module
	PublishSwap(context.Context, model.Swap) error
	PublishLiquidityPool(context.Context, model.LiquidityPool) error

	// graph module
	PublishCyberlink(context.Context, model.Cyberlink) error
	PublishCyberlinkMessage(context.Context, model.CyberlinkMessage) error
	PublishParticle(context.Context, model.Particle) error

	// bandwidth module
	PublishBandwidthParams(context.Context, model.BandwidthParams) error

	// dmn module
	PublishDMNParams(context.Context, model.DMNParams) error

	// grid module
	PublishGridParams(context.Context, model.GridParams) error
	PublishRoute(context.Context, model.Route) error
	PublishCreateRouteMessage(context.Context, model.CreateRouteMessage) error
	PublishEditRouteMessage(context.Context, model.EditRouteMessage) error
	PublishEditRouteNameMessage(context.Context, model.EditRouteNameMessage) error
	PublishDeleteRouteMessage(context.Context, model.DeleteRouteMessage) error

	// rank module
	PublishRankParams(context.Context, model.RankParams) error

	// resources module
	PublishInvestmintMessage(context.Context, model.InvestmintMessage) error

	// raw
	PublishRawBlock(ctx context.Context, b interface{}) error
	PublishRawTransaction(ctx context.Context, tx interface{}) error
	PublishRawBlockResults(ctx context.Context, br interface{}) error
	PublishRawGenesis(ctx context.Context, g interface{}) error
}
