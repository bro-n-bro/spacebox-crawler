package broker

var (
	Account                    Topic = newTopic("account")
	AccountBalance             Topic = newTopic("account_balance")
	AnnualProvision            Topic = newTopic("annual_provision")
	Block                      Topic = newTopic("block")
	CommunityPool              Topic = newTopic("community_pool")
	DistributionParams         Topic = newTopic("distribution_params")
	DelegationReward           Topic = newTopic("delegation_reward")
	DelegationRewardMessage    Topic = newTopic("delegation_reward_message")
	Delegation                 Topic = newTopic("delegation")
	DelegationMessage          Topic = newTopic("delegation_message")
	GovParams                  Topic = newTopic("gov_params")
	Message                    Topic = newTopic("message")
	MintParams                 Topic = newTopic("mint_params")
	MultiSendMessage           Topic = newTopic("multisend_message")
	Proposal                   Topic = newTopic("proposal")
	ProposalVoteMessage        Topic = newTopic("proposal_vote_message")
	ProposalTallyResult        Topic = newTopic("proposal_tally_result")
	ProposalDeposit            Topic = newTopic("proposal_deposit")
	ProposalDepositMessage     Topic = newTopic("proposal_deposit_message")
	ProposerReward             Topic = newTopic("proposer_reward")
	Redelegation               Topic = newTopic("redelegation")
	RedelegationMessage        Topic = newTopic("redelegation_message")
	SendMessage                Topic = newTopic("send_message")
	SetWithdrawAddressMessage  Topic = newTopic("set_withdraw_address_message")
	StakingParams              Topic = newTopic("staking_params")
	StakingPool                Topic = newTopic("staking_pool")
	Supply                     Topic = newTopic("supply")
	Transaction                Topic = newTopic("transaction")
	UnbondingDelegation        Topic = newTopic("unbonding_delegation")
	UnbondingDelegationMessage Topic = newTopic("unbonding_delegation_message")
	Validator                  Topic = newTopic("validator")
	ValidatorInfo              Topic = newTopic("validator_info")
	ValidatorStatus            Topic = newTopic("validator_status")
	ValidatorDescription       Topic = newTopic("validator_description")
	ValidatorCommission        Topic = newTopic("validator_commission")

	authTopics = Topics{Account}

	bankTopics = Topics{Supply, AccountBalance, SendMessage, MultiSendMessage}

	distributionTopics = Topics{DistributionParams, CommunityPool, /* TODO: validatorCommission, DelegationReward, */
		DelegationRewardMessage, SetWithdrawAddressMessage, ProposerReward}

	govTopics = Topics{GovParams, Proposal, ProposalDepositMessage, ProposalTallyResult, ProposalVoteMessage}

	mintTopics = Topics{MintParams, AnnualProvision}

	stakingTopics = Topics{Validator, ValidatorStatus, ValidatorInfo, ValidatorDescription, StakingParams,
		StakingPool, Redelegation, RedelegationMessage, UnbondingDelegation, UnbondingDelegationMessage,
		Delegation, DelegationMessage,
	}

	coreTopics = Topics{Block, Transaction, Message}
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
