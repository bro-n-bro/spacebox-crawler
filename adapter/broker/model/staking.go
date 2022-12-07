package model

import (
	"time"
)

type (
	UnbondingDelegation struct {
		CompletionTimestamp time.Time `json:"completion_timestamp"`
		Coin                Coin      `json:"coin"`
		DelegatorAddress    string    `json:"delegator_address"`
		ValidatorOperAddr   string    `json:"validator_oper_addr"`
		TxHash              string    `json:"tx_hash"`
		Height              int64     `json:"height"`
	}

	UnbondingDelegationMessage struct {
		UnbondingDelegation
		TxHash string `json:"tx_hash"`
	}

	Delegation struct {
		OperatorAddress  string `json:"operator_address"`
		DelegatorAddress string `json:"delegator_address"`
		Coin             Coin   `json:"coin"`
		Height           int64  `json:"height"`
	}
	DelegationMessage struct {
		Delegation
		TxHash string `json:"tx_hash"`
	}

	DelegationRewardMessage struct {
		Coins            Coins  `json:"coins"`
		Height           int64  `json:"height"`
		DelegatorAddress string `json:"delegator_address"`
		ValidatorAddress string `json:"validator_address"`
		TxHash           string `json:"tx_hash"`
	}
)
