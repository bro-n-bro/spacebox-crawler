package broker

type (
	Topic  *string
	Topics []Topic
)

var (
	Account                    Topic = newTopic("account")
	AccountBalance             Topic = newTopic("account_balance")
	AnnualProvision            Topic = newTopic("annual_provision")
	Block                      Topic = newTopic("block")
	CommunityPool              Topic = newTopic("community_pool")
	DistributionParams         Topic = newTopic("distribution_params")
	DelegationRewardMessage    Topic = newTopic("delegation_reward_message")
	Delegation                 Topic = newTopic("delegation")
	DelegationMessage          Topic = newTopic("delegation_message")
	GovParams                  Topic = newTopic("gov_params")
	Message                    Topic = newTopic("message")
	MintParams                 Topic = newTopic("mint_params")
	MultiSendMessage           Topic = newTopic("multisend_message")
	ProposalVoteMessage        Topic = newTopic("proposal_vote_message")
	ProposalTallyResult        Topic = newTopic("proposal_tally_result")
	Redelegation               Topic = newTopic("redelegation")
	RedelegationMessage        Topic = newTopic("redelegation_message")
	SendMessage                Topic = newTopic("send_message")
	StakingParams              Topic = newTopic("staking_params")
	StakingPool                Topic = newTopic("staking_pool")
	Supply                     Topic = newTopic("supply")
	Transaction                Topic = newTopic("tx")
	UnbondingDelegationMessage Topic = newTopic("unbonding_delegation_message")
	UnbondingDelegation        Topic = newTopic("unbonding_delegation")
	ValidatorInfo              Topic = newTopic("validator_info")
	ValidatorStatus            Topic = newTopic("validator_status")
	Validator                  Topic = newTopic("validator")

	authTopics = Topics{Account}

	bankTopics = Topics{Supply, AccountBalance, SendMessage, MultiSendMessage}

	distributionTopics = Topics{DistributionParams, CommunityPool, /* TODO: validatorCommission, DelegationRewardMessage */
		DelegationRewardMessage}

	govTopics = Topics{GovParams /*TODO: Proposal, ProposalDepositMessage */, ProposalTallyResult, ProposalVoteMessage,
		MultiSendMessage}

	mintTopics = Topics{MintParams, AnnualProvision}

	stakingTopics = Topics{Validator, ValidatorStatus, ValidatorInfo, StakingParams, StakingPool, Redelegation,
		RedelegationMessage, UnbondingDelegation, UnbondingDelegationMessage, Delegation, DelegationMessage}

	coreTopics = Topics{Block, Transaction, Message}
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
