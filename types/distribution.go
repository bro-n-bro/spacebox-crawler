package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
)

type (
	// DistributionParams represents the parameters of the x/distribution module
	DistributionParams struct {
		distrtypes.Params
		Height int64
	}

	// ValidatorCommissionAmount represents the commission amount for a specific validator
	ValidatorCommissionAmount struct {
		ValidatorOperAddr         string
		ValidatorSelfDelegateAddr string
		Amount                    []sdk.DecCoin
		Height                    int64
	}
	// DelegatorRewardAmount contains the data of a delegator commission amount
	DelegatorRewardAmount struct {
		ValidatorOperAddr string
		DelegatorAddress  string
		WithdrawAddress   string
		Amount            []sdk.DecCoin
		Height            int64
	}
	// DelegationRewardMessage contains Coins for DelegatorAddress and ValidatorAddress
	// coming from MsgWithdrawDelegatorReward tx message
	DelegationRewardMessage struct {
		DelegatorAddress string
		ValidatorAddress string
		Coins            Coins
	}
)

// NewDistributionParams allows to build a new DistributionParams instance
func NewDistributionParams(params distrtypes.Params, height int64) DistributionParams {
	return DistributionParams{
		Params: params,
		Height: height,
	}
}

// NewValidatorCommissionAmount allows to build a new ValidatorCommissionAmount instance
func NewValidatorCommissionAmount(
	valOperAddr, valSelfDelegateAddress string, amount sdk.DecCoins, height int64,
) ValidatorCommissionAmount {
	return ValidatorCommissionAmount{
		ValidatorOperAddr:         valOperAddr,
		ValidatorSelfDelegateAddr: valSelfDelegateAddress,
		Amount:                    amount,
		Height:                    height,
	}
}

// NewDelegatorRewardAmount allows to build a new DelegatorRewardAmount instance
func NewDelegatorRewardAmount(
	delegator, valOperAddr, withdrawAddress string, amount sdk.DecCoins, height int64,
) DelegatorRewardAmount {
	return DelegatorRewardAmount{
		ValidatorOperAddr: valOperAddr,
		DelegatorAddress:  delegator,
		WithdrawAddress:   withdrawAddress,
		Amount:            amount,
		Height:            height,
	}
}

// NewDelegationRewardMessage returns the new instance of DelegationRewardMessage
func NewDelegationRewardMessage(delAddr, valAddr string, coins Coins) DelegationRewardMessage {
	return DelegationRewardMessage{
		DelegatorAddress: delAddr,
		ValidatorAddress: valAddr,
		Coins:            coins,
	}
}
