package broker

type Topic *string

var (
	Account                    Topic = newTopic("account")
	AccountBalance             Topic = newTopic("account_balance")
	BlockTopic                 Topic = newTopic("block")
	CommunityPool              Topic = newTopic("community_pool")
	DistributionParams         Topic = newTopic("distribution_params")
	DelegationRewardMessage    Topic = newTopic("delegation_reward_message")
	Delegation                 Topic = newTopic("delegation")
	DelegationMessage          Topic = newTopic("delegation_message")
	GovParams                  Topic = newTopic("gov_params")
	MessageTopic               Topic = newTopic("message")
	MintParams                 Topic = newTopic("mint_params")
	MultiSendMessage           Topic = newTopic("multisend_message")
	ProposalVoteMessage        Topic = newTopic("proposal_vote_message")
	ProposalTallyResult        Topic = newTopic("proposal_tally_result")
	Redelegation               Topic = newTopic("redelegation")
	RedelegationMessage        Topic = newTopic("redelegation_message")
	StakingParams              Topic = newTopic("staking_params")
	SendMessage                Topic = newTopic("send_message")
	SupplyTopic                Topic = newTopic("supply")
	Transaction                Topic = newTopic("tx")
	UnbondingDelegationMessage Topic = newTopic("unbonding_delegation_message")
	UnbondingDelegation        Topic = newTopic("unbonding_delegation")
	ValidatorInfo              Topic = newTopic("validator_info")
	ValidatorStatus            Topic = newTopic("validator_status")
	Validator                  Topic = newTopic("validator")
)

func newTopic(t string) *string { return &t }
