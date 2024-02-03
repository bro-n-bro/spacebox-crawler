package broker

var (
	Account                            Topic = newTopic("account")
	AccountBalance                     Topic = newTopic("account_balance")
	AcknowledgementMessage             Topic = newTopic("acknowledgement_message")
	AnnualProvision                    Topic = newTopic("annual_provision")
	AuthzGrant                         Topic = newTopic("authz_grant")
	BandwidthParams                    Topic = newTopic("bandwidth_params")
	Block                              Topic = newTopic("block")
	CancelUnbondingDelegationMessage   Topic = newTopic("cancel_unbonding_delegation_message")
	CommunityPool                      Topic = newTopic("community_pool")
	CreateRouteMessage                 Topic = newTopic("create_route_message")
	CreateValidatorMessage             Topic = newTopic("create_validator_message")
	Cyberlink                          Topic = newTopic("cyberlink")
	CyberlinkMessage                   Topic = newTopic("cyberlink_message")
	DMNParams                          Topic = newTopic("dmn_params")
	Delegation                         Topic = newTopic("delegation")
	DelegationMessage                  Topic = newTopic("delegation_message")
	DelegationReward                   Topic = newTopic("delegation_reward")
	DelegationRewardMessage            Topic = newTopic("delegation_reward_message")
	DeleteRouteMessage                 Topic = newTopic("delete_route_message")
	DenomTrace                         Topic = newTopic("denom_trace")
	DistributionCommission             Topic = newTopic("distribution_commission")
	DistributionParams                 Topic = newTopic("distribution_params")
	DistributionReward                 Topic = newTopic("distribution_reward")
	EditRouteMessage                   Topic = newTopic("edit_route_message")
	EditRouteNameMessage               Topic = newTopic("edit_route_name_message")
	EditValidatorMessage               Topic = newTopic("edit_validator_message")
	ExecMessage                        Topic = newTopic("exec_message")
	FeeAllowance                       Topic = newTopic("fee_allowance")
	GovParams                          Topic = newTopic("gov_params")
	GrantAllowanceMessage              Topic = newTopic("grant_allowance_message")
	GrantMessage                       Topic = newTopic("grant_message")
	GridParams                         Topic = newTopic("grid_params")
	HandleValidatorSignature           Topic = newTopic("handle_validator_signature")
	InvestmintMessage                  Topic = newTopic("investmint_message")
	LiquidityPool                      Topic = newTopic("liquidity_pool")
	Message                            Topic = newTopic("message")
	MintParams                         Topic = newTopic("mint_params")
	MultiSendMessage                   Topic = newTopic("multisend_message")
	Particle                           Topic = newTopic("particle")
	Proposal                           Topic = newTopic("proposal")
	ProposalDeposit                    Topic = newTopic("proposal_deposit")
	ProposalDepositMessage             Topic = newTopic("proposal_deposit_message")
	ProposalTallyResult                Topic = newTopic("proposal_tally_result")
	ProposalVoteMessage                Topic = newTopic("proposal_vote_message")
	ProposerReward                     Topic = newTopic("proposer_reward")
	RankParams                         Topic = newTopic("rank_params")
	RawBlock                           Topic = newTopic("raw_block")
	RawBlockResults                    Topic = newTopic("raw_block_results")
	RawGenesis                         Topic = newTopic("raw_genesis")
	RawTransaction                     Topic = newTopic("raw_transaction")
	ReceivePacketMessage               Topic = newTopic("receive_packet_message")
	Redelegation                       Topic = newTopic("redelegation")
	RedelegationMessage                Topic = newTopic("redelegation_message")
	RevokeAllowanceMessage             Topic = newTopic("revoke_allowance_message")
	RevokeMessage                      Topic = newTopic("revoke_message")
	Route                              Topic = newTopic("route")
	SendMessage                        Topic = newTopic("send_message")
	SetWithdrawAddressMessage          Topic = newTopic("set_withdraw_address_message")
	SlashingParams                     Topic = newTopic("slashing_params")
	StakingParams                      Topic = newTopic("staking_params")
	StakingPool                        Topic = newTopic("staking_pool")
	SubmitProposalMessage              Topic = newTopic("submit_proposal_message")
	Supply                             Topic = newTopic("supply")
	Swap                               Topic = newTopic("swap")
	Transaction                        Topic = newTopic("transaction")
	TransferMessage                    Topic = newTopic("transfer_message")
	UnbondingDelegation                Topic = newTopic("unbonding_delegation")
	UnbondingDelegationMessage         Topic = newTopic("unbonding_delegation_message")
	UnjailMessage                      Topic = newTopic("unjail_message")
	Validator                          Topic = newTopic("validator")
	ValidatorCommission                Topic = newTopic("validator_commission")
	ValidatorDescription               Topic = newTopic("validator_description")
	ValidatorInfo                      Topic = newTopic("validator_info")
	ValidatorPreCommit                 Topic = newTopic("validator_precommit")
	ValidatorStatus                    Topic = newTopic("validator_status")
	ValidatorVotingPower               Topic = newTopic("validator_voting_power")
	VoteWeightedMessage                Topic = newTopic("vote_weighted_message")
	WithdrawValidatorCommissionMessage Topic = newTopic("withdraw_validator_commission_message")

	authTopics = Topics{Account}

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

	coreTopics = Topics{Block, Transaction, Message, ValidatorPreCommit}

	authzTopics = Topics{AuthzGrant, GrantMessage, RevokeMessage, ExecMessage}

	feegrantTopics = Topics{FeeAllowance, GrantAllowanceMessage, RevokeAllowanceMessage}

	slashingTopics = Topics{UnjailMessage, HandleValidatorSignature, SlashingParams}

	ibcTopics = Topics{TransferMessage, AcknowledgementMessage, ReceivePacketMessage, DenomTrace}

	liquidityTopics = Topics{Swap, LiquidityPool}

	graphTopics = Topics{Cyberlink, CyberlinkMessage, Particle}

	bandwidthTopics = Topics{BandwidthParams}

	dmnTopics = Topics{DMNParams}

	gridTopics = Topics{GridParams, CreateRouteMessage, EditRouteMessage, EditRouteNameMessage, DeleteRouteMessage, Route}

	rankTopics = Topics{RankParams}

	resourcesTopics = Topics{InvestmintMessage}

	wasmTopics = Topics{Cyberlink, Particle}

	rawTopics = Topics{RawBlock, RawTransaction, RawBlockResults, RawGenesis}
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
