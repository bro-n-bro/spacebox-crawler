package types

import (
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// StakingValidator represents a single validator.
// This is defined as an interface so that we can use the SDK types
// as well as database types properly.
type (
	StakingValidator interface {
		GetConsAddr() string
		GetConsPubKey() string
		GetOperator() string
		GetSelfDelegateAddress() string
		GetMaxChangeRate() *sdk.Dec
		GetMaxRate() *sdk.Dec
		GetHeight() int64
		GetMinSelfDelegation() *sdkmath.Int
	}

	// validator allows to easily implement the Validator interface
	stakingValidator struct {
		MinSelfDelegation   *sdkmath.Int
		MaxChangeRate       *sdk.Dec
		MaxRate             *sdk.Dec
		ConsPubKey          string
		OperatorAddr        string
		SelfDelegateAddress string
		ConsensusAddr       string
		Height              int64
	}
)

// NewStakingValidator allows to build a new Validator implementation having the given data
func NewStakingValidator(consAddr, opAddr, consPubKey, selfDelegateAddress string,
	maxChangeRate, maxRate *sdk.Dec, height int64) StakingValidator {

	return stakingValidator{
		ConsensusAddr:       consAddr,
		ConsPubKey:          consPubKey,
		OperatorAddr:        opAddr,
		SelfDelegateAddress: selfDelegateAddress,
		MaxChangeRate:       maxChangeRate,
		MaxRate:             maxRate,
		Height:              height,
	}
}

// GetConsAddr implements the Validator interface
func (v stakingValidator) GetConsAddr() string {
	return v.ConsensusAddr
}

// GetConsPubKey implements the Validator interface
func (v stakingValidator) GetConsPubKey() string {
	return v.ConsPubKey
}

func (v stakingValidator) GetOperator() string {
	return v.OperatorAddr
}

func (v stakingValidator) GetSelfDelegateAddress() string {
	return v.SelfDelegateAddress
}

func (v stakingValidator) GetMaxChangeRate() *sdk.Dec {
	return v.MaxChangeRate
}

func (v stakingValidator) GetMaxRate() *sdk.Dec {
	return v.MaxRate
}

func (v stakingValidator) GetHeight() int64 {
	return v.Height
}

func (v stakingValidator) GetMinSelfDelegation() *sdkmath.Int {
	return v.MinSelfDelegation
}
