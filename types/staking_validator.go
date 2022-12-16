package types

import (
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

// Validator represents a single validator.
// This is defined as an interface so that we can use the SDK types
// as well as database types properly.
type StakingValidator interface {
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
type stakingValidator struct {
	ConsensusAddr       string
	ConsPubKey          string
	OperatorAddr        string
	SelfDelegateAddress string
	MaxChangeRate       *sdk.Dec
	MaxRate             *sdk.Dec
	MinSelfDelegation   *sdkmath.Int
	Height              int64
}

// NewStakingValidator allows to build a new Validator implementation having the given data
func NewStakingValidator(
	consAddr, opAddr, consPubKey, selfDelegateAddress string, maxChangeRate, maxRate *sdk.Dec,
	height int64) StakingValidator {
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

// --------------------------------------------------------------------------------------------------------------------

// ValidatorDescription contains the description of a validator
// and timestamp do the description get changed
type ValidatorDescription struct {
	OperatorAddress string
	Description     stakingtypes.Description
	AvatarURL       string
	Height          int64
}

// NewValidatorDescription return a new ValidatorDescription object
func NewValidatorDescription(
	opAddr string, description stakingtypes.Description, avatarURL string, height int64,
) ValidatorDescription {
	return ValidatorDescription{
		OperatorAddress: opAddr,
		Description:     description,
		AvatarURL:       avatarURL,
		Height:          height,
	}
}

// ----------------------------------------------------------------------------------------------------------

// ValidatorCommission contains the data of a validator commission at a given height
type ValidatorCommission struct {
	ValAddress    string
	Commission    *sdk.Dec
	MaxChangeRate *sdk.Dec
	MaxRate       *sdk.Dec
	Height        int64
}

// NewValidatorCommission return a new validator commission instance
func NewValidatorCommission(valAddress string, rate, maxChangeRate, maxRate *sdk.Dec, height int64) ValidatorCommission {
	return ValidatorCommission{
		ValAddress:    valAddress,
		Commission:    rate,
		MaxChangeRate: maxChangeRate,
		MaxRate:       maxRate,
		Height:        height,
	}
}

//--------------------------------------------

// ValidatorVotingPower represents the voting power of a validator at a specific block height
type ValidatorVotingPower struct {
	ConsensusAddress string
	VotingPower      int64
	Height           int64
}

// NewValidatorVotingPower creates a new ValidatorVotingPower
func NewValidatorVotingPower(address string, votingPower int64, height int64) ValidatorVotingPower {
	return ValidatorVotingPower{
		ConsensusAddress: address,
		VotingPower:      votingPower,
		Height:           height,
	}
}

//--------------------------------------------------------

// ValidatorStatus represents the current state for the specified validator at the specific height
type ValidatorStatus struct {
	ConsensusAddress string
	ConsensusPubKey  string
	Status           int
	Jailed           bool
	Height           int64
}

// NewValidatorStatus creates a new ValidatorVotingPower
func NewValidatorStatus(valConsAddr, pubKey string, status int, jailed bool, height int64) ValidatorStatus {
	return ValidatorStatus{
		ConsensusAddress: valConsAddr,
		ConsensusPubKey:  pubKey,
		Status:           status,
		Jailed:           jailed,
		Height:           height,
	}
}
