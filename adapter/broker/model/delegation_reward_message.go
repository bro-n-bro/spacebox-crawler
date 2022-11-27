package model

type DelegationRewardMessage struct {
	Coins            Coins  `json:"coins"`
	DelegatorAddress string `json:"delegator_address"`
	ValidatorAddress string `json:"validator_address"`
}
