package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

type (
	StakingValidator interface {
		GetConsAddr() string
		GetConsPubKey() string
		GetOperator() string
		GetSelfDelegateAddress() string
		GetMaxChangeRate() *sdk.Dec
		GetMaxRate() *sdk.Dec
		GetHeight() int64
		GetMinSelfDelegation() int64
		GetDescription() stakingtypes.Description
	}

	// validator allows to easily implement the Validator interface
	stakingValidator struct {
		MaxChangeRate       *sdk.Dec
		MaxRate             *sdk.Dec
		Description         stakingtypes.Description
		ConsPubKey          string
		OperatorAddr        string
		SelfDelegateAddress string
		ConsensusAddr       string
		Height              int64
		MinSelfDelegation   int64
	}
)

// NewStakingValidator allows to build a new Validator implementation having the given data
func NewStakingValidator(
	consensusAddr,
	operatorAddr,
	consensusPubKey,
	selfDelegateAddress string,
	maxChangeRate,
	maxRate *sdk.Dec,
	description stakingtypes.Description,
	height,
	minSelfDelegation int64,
) StakingValidator {

	return stakingValidator{
		ConsensusAddr:       consensusAddr,
		ConsPubKey:          consensusPubKey,
		OperatorAddr:        operatorAddr,
		SelfDelegateAddress: selfDelegateAddress,
		MaxChangeRate:       maxChangeRate,
		MaxRate:             maxRate,
		Description:         description,
		Height:              height,
		MinSelfDelegation:   minSelfDelegation,
	}
}

func (v stakingValidator) GetSelfDelegateAddress() string           { return v.SelfDelegateAddress }
func (v stakingValidator) GetMinSelfDelegation() int64              { return v.MinSelfDelegation }
func (v stakingValidator) GetMaxChangeRate() *sdk.Dec               { return v.MaxChangeRate }
func (v stakingValidator) GetDescription() stakingtypes.Description { return v.Description }
func (v stakingValidator) GetConsPubKey() string                    { return v.ConsPubKey }
func (v stakingValidator) GetOperator() string                      { return v.OperatorAddr }
func (v stakingValidator) GetConsAddr() string                      { return v.ConsensusAddr }
func (v stakingValidator) GetMaxRate() *sdk.Dec                     { return v.MaxRate }
func (v stakingValidator) GetHeight() int64                         { return v.Height }
