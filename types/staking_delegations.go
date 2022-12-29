package types

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type (
	// Delegation represents a single delegation made from a delegator
	// to a specific validator at a specific height (and timestamp)
	// containing a given amount of tokens
	Delegation struct {
		DelegatorAddress  string
		ValidatorOperAddr string
		Coin              Coin
		Height            int64
	}

	DelegationMessage struct {
		TxHash string
		Delegation
	}

	// UnbondingDelegation represents a single unbonding delegation
	UnbondingDelegation struct {
		CompletionTimestamp time.Time
		DelegatorAddress    string
		ValidatorOperAddr   string
		Coin                Coin
		Height              int64
	}

	// UnbondingDelegationMessage
	UnbondingDelegationMessage struct {
		TxHash string
		Coin   Coin
		UnbondingDelegation
	}

	// Redelegation represents a single re-delegations
	Redelegation struct {
		CompletionTime   time.Time
		DelegatorAddress string
		SrcValidator     string
		DstValidator     string
		Coin             Coin
		Height           int64
	}

	// RedelegationMessage
	RedelegationMessage struct {
		TxHash string
		Coin   Coin
		Redelegation
	}
)

// NewDelegation creates a new Delegation instance containing
// the given data
func NewDelegation(delegator, validatorOperAddr string, amount sdk.Coin, height int64) Delegation {
	return Delegation{
		DelegatorAddress:  delegator,
		ValidatorOperAddr: validatorOperAddr,
		Coin:              NewCoinFromCdk(amount),
		Height:            height,
	}
}

// NewDelegationMessage creates a new DelegationMessage instance
func NewDelegationMessage(delegator, validatorOperAddr, txHash string, amount sdk.Coin, height int64) DelegationMessage {
	return DelegationMessage{
		Delegation: NewDelegation(delegator, validatorOperAddr, amount, height),
		TxHash:     txHash,
	}
}

// NewUnbondingDelegation allows to create a new UnbondingDelegation instance
func NewUnbondingDelegation(delegator, validatorOperAddr string, coin sdk.Coin, completionTimestamp time.Time,
	height int64) UnbondingDelegation {
	return UnbondingDelegation{
		DelegatorAddress:    delegator,
		ValidatorOperAddr:   validatorOperAddr,
		Coin:                NewCoinFromCdk(coin),
		CompletionTimestamp: completionTimestamp,
		Height:              height,
	}
}

// NewUnbondingDelegationMessage
func NewUnbondingDelegationMessage(delegator, validatorOperAddr, txHash string, coin Coin, completionTimestamp time.Time,
	height int64) UnbondingDelegationMessage {
	return UnbondingDelegationMessage{
		UnbondingDelegation: UnbondingDelegation{
			DelegatorAddress:    delegator,
			ValidatorOperAddr:   validatorOperAddr,
			CompletionTimestamp: completionTimestamp,
			Height:              height,
		},
		Coin:   coin,
		TxHash: txHash,
	}
}

// Equal returns true iff u and v contain the same data
func (u UnbondingDelegation) Equal(v UnbondingDelegation) bool {
	return u.DelegatorAddress == v.DelegatorAddress &&
		u.ValidatorOperAddr == v.ValidatorOperAddr &&
		u.Coin.IsEqual(v.Coin) &&
		u.CompletionTimestamp.Equal(v.CompletionTimestamp) &&
		u.Height == v.Height
}

// NewRedelegation build a new Redelegation object
func NewRedelegation(delegator, srcValidator, dstValidator string, amount sdk.Coin, completionTime time.Time,
	height int64) Redelegation {
	return Redelegation{
		DelegatorAddress: delegator,
		SrcValidator:     srcValidator,
		DstValidator:     dstValidator,
		Coin:             NewCoinFromCdk(amount),
		CompletionTime:   completionTime,
		Height:           height,
	}
}

// Equal returns true iff r and s contain the same data
func (r Redelegation) Equal(s Redelegation) bool {
	return r.DelegatorAddress == s.DelegatorAddress &&
		r.SrcValidator == s.SrcValidator &&
		r.DstValidator == s.DstValidator &&
		r.Coin.IsEqual(s.Coin) &&
		r.CompletionTime.Equal(s.CompletionTime) &&
		r.Height == s.Height
}

func NewRedelegationMessage(delegator, srcValidator, dstValidator, txHash string, coin Coin,
	completionTime time.Time, height int64) RedelegationMessage {
	return RedelegationMessage{
		Redelegation: Redelegation{
			DelegatorAddress: delegator,
			SrcValidator:     srcValidator,
			DstValidator:     dstValidator,
			CompletionTime:   completionTime,
			Height:           height,
		},
		Coin:   coin,
		TxHash: txHash,
	}
}
