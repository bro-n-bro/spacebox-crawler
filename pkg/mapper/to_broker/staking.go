package to_broker

import (
	"github.com/hexy-dev/spacebox/broker/model"

	"github.com/hexy-dev/spacebox-crawler/types"
)

func (tb ToBroker) MapUnbondingDelegation(ud types.UnbondingDelegation) model.UnbondingDelegation {
	return model.UnbondingDelegation{
		CompletionTimestamp: ud.CompletionTimestamp,
		Coin:                tb.MapCoin(ud.Coin),
		DelegatorAddress:    ud.DelegatorAddress,
		ValidatorAddress:    ud.ValidatorOperAddr,
		Height:              ud.Height,
	}
}

func (tb ToBroker) MapUnbondingDelegationMessage(udm types.UnbondingDelegationMessage) model.UnbondingDelegationMessage {
	return model.UnbondingDelegationMessage{
		UnbondingDelegation: tb.MapUnbondingDelegation(udm.UnbondingDelegation),
		TxHash:              udm.TxHash,
	}
}
func (tb ToBroker) MapStakingParams(sp types.StakingParams) model.StakingParams {
	var commissionRate float64
	if !sp.MinCommissionRate.IsNil() {
		commissionRate = sp.MinCommissionRate.MustFloat64()
	}
	return model.NewStakingParams(sp.Height, sp.MaxValidators, sp.MaxEntries, sp.HistoricalEntries,
		sp.BondDenom, commissionRate, sp.UnbondingTime)
}

func (tb ToBroker) MapStakingPool(sp *types.Pool) model.StakingPool {
	pool := model.StakingPool{
		Height: sp.Height,
	}
	if !sp.BondedTokens.IsNil() {
		pool.BondedTokens = sp.BondedTokens.Int64()
	}

	if !sp.NotBondedTokens.IsNil() {
		pool.NotBondedTokens = sp.NotBondedTokens.Int64()
	}
	return pool
}

func (tb ToBroker) MapDelegation(d types.Delegation) model.Delegation {
	return model.Delegation{
		OperatorAddress:  d.ValidatorOperAddr,
		DelegatorAddress: d.DelegatorAddress,
		Coin:             tb.MapCoin(d.Coin),
		Height:           d.Height,
	}
}

func (tb ToBroker) MapDelegationMessage(d types.DelegationMessage) model.DelegationMessage {
	return model.DelegationMessage{
		Delegation: tb.MapDelegation(d.Delegation),
		TxHash:     d.TxHash,
	}
}
