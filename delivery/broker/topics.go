package broker

var (
	Account                            Topic = newTopic("account")
	AuthzGrant                         Topic = newTopic("authz_grant")
	AccountBalance                     Topic = newTopic("account_balance")
	AnnualProvision                    Topic = newTopic("annual_provision")
	Block                              Topic = newTopic("block")
	CancelUnbondingDelegationMessage   Topic = newTopic("cancel_unbonding_delegation_message")
	CommunityPool                      Topic = newTopic("community_pool")
	CreateValidatorMessage             Topic = newTopic("create_validator_message")
	DistributionCommission             Topic = newTopic("distribution_commission")
	DistributionReward                 Topic = newTopic("distribution_reward")
	DistributionParams                 Topic = newTopic("distribution_params")
	DelegationReward                   Topic = newTopic("delegation_reward")
	DelegationRewardMessage            Topic = newTopic("delegation_reward_message")
	Delegation                         Topic = newTopic("delegation")
	DelegationMessage                  Topic = newTopic("delegation_message")
	EditValidatorMessage               Topic = newTopic("edit_validator_message")
	ExecMessage                        Topic = newTopic("exec_message")
	FeeAllowance                       Topic = newTopic("fee_allowance")
	GovParams                          Topic = newTopic("gov_params")
	GrantMessage                       Topic = newTopic("grant_message")
	GrantAllowanceMessage              Topic = newTopic("grant_allowance_message")
	Message                            Topic = newTopic("message")
	MintParams                         Topic = newTopic("mint_params")
	MultiSendMessage                   Topic = newTopic("multisend_message")
	Proposal                           Topic = newTopic("proposal")
	ProposalVoteMessage                Topic = newTopic("proposal_vote_message")
	ProposalTallyResult                Topic = newTopic("proposal_tally_result")
	ProposalDeposit                    Topic = newTopic("proposal_deposit")
	ProposalDepositMessage             Topic = newTopic("proposal_deposit_message")
	ProposerReward                     Topic = newTopic("proposer_reward")
	RevokeAllowanceMessage             Topic = newTopic("revoke_allowance_message")
	Redelegation                       Topic = newTopic("redelegation")
	RedelegationMessage                Topic = newTopic("redelegation_message")
	RevokeMessage                      Topic = newTopic("revoke_message")
	SendMessage                        Topic = newTopic("send_message")
	SetWithdrawAddressMessage          Topic = newTopic("set_withdraw_address_message")
	StakingParams                      Topic = newTopic("staking_params")
	StakingPool                        Topic = newTopic("staking_pool")
	SubmitProposalMessage              Topic = newTopic("submit_proposal_message")
	Supply                             Topic = newTopic("supply")
	Transaction                        Topic = newTopic("transaction")
	UnbondingDelegation                Topic = newTopic("unbonding_delegation")
	UnbondingDelegationMessage         Topic = newTopic("unbonding_delegation_message")
	Validator                          Topic = newTopic("validator")
	ValidatorInfo                      Topic = newTopic("validator_info")
	ValidatorStatus                    Topic = newTopic("validator_status")
	ValidatorDescription               Topic = newTopic("validator_description")
	ValidatorCommission                Topic = newTopic("validator_commission")
	VoteWeightedMessage                Topic = newTopic("vote_weighted_message")
	WithdrawValidatorCommissionMessage Topic = newTopic("withdraw_validator_commission_message")
	authTopics                               = Topics{Account}

	bankTopics = Topics{Supply, AccountBalance, SendMessage, MultiSendMessage}

	distributionTopics = Topics{DistributionCommission, DistributionParams, CommunityPool,
		DelegationRewardMessage, SetWithdrawAddressMessage, ProposerReward, DistributionReward,
		WithdrawValidatorCommissionMessage, /* TODO: validatorCommission, DelegationReward, */
	}

	govTopics = Topics{GovParams, Proposal, ProposalDepositMessage, ProposalTallyResult, ProposalVoteMessage,
		VoteWeightedMessage, SubmitProposalMessage}

	mintTopics = Topics{MintParams, AnnualProvision}

	stakingTopics = Topics{Validator, ValidatorStatus, ValidatorInfo, ValidatorDescription, StakingParams,
		StakingPool, Redelegation, RedelegationMessage, UnbondingDelegation, UnbondingDelegationMessage,
		Delegation, DelegationMessage, CreateValidatorMessage, EditValidatorMessage, CancelUnbondingDelegationMessage,
	}

	coreTopics = Topics{Block, Transaction, Message}

	authzTopics = Topics{AuthzGrant, GrantMessage, RevokeMessage, ExecMessage}

	feegrantTopics = Topics{FeeAllowance, GrantAllowanceMessage, RevokeAllowanceMessage}
)

type (
	Topic  *string
	Topics []Topic
)

func newTopic(t string) *string { return &t }

func (ts Topics) ToStringSlice() []string {
	res := make([]string, len(ts))

	for i, t := range ts {
		if t == nil {
			panic("topic is nil")
		}

		res[i] = *t
	}

	return res
}
