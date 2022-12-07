package broker

type Topic *string

var (
	AccountTopic                 Topic = newTopic("account")
	AccountBalance               Topic = newTopic("account_balance")
	BlockTopic                   Topic = newTopic("block")
	DistributionParamsTopic      Topic = newTopic("distribution_params")
	DelegationRewardMessageTopic Topic = newTopic("delegation_reward_message")
	Delegation                   Topic = newTopic("delegation")
	DelegationMessage            Topic = newTopic("delegation_message")
	GovParams                    Topic = newTopic("gov_params")
	MessageTopic                 Topic = newTopic("message")
	MintParams                   Topic = newTopic("mint_params")
	MultiSendMessageTopic        Topic = newTopic("multisend_message")
	ProposalVoteMessageTopic     Topic = newTopic("proposal_vote_message")
	ProposalTallyResult          Topic = newTopic("proposal_tally_result")
	RedelegationMessageTopic     Topic = newTopic("redelegation_message")
	StakingParams                Topic = newTopic("staking_params")
	SendMessageTopic             Topic = newTopic("send_message")
	SupplyTopic                  Topic = newTopic("supply")
	TransactionTopic             Topic = newTopic("tx")
	UnbondingDelegationMessage   Topic = newTopic("unbonding_delegation_message")
	UnbondingDelegation          Topic = newTopic("unbonding_delegation")
	ValidatorInfo                Topic = newTopic("validator_info")
	Validator                    Topic = newTopic("validator")
)

func newTopic(t string) *string { return &t }
