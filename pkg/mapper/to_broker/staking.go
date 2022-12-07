package to_broker

import (
	"bro-n-bro-osmosis/adapter/broker/model"
	"bro-n-bro-osmosis/types"
)

func (tb ToBroker) MapUnbondingDelegation(ud types.UnbondingDelegation) model.UnbondingDelegation {
	return model.UnbondingDelegation{
		CompletionTimestamp: ud.CompletionTimestamp,
		Coin:                tb.MapCoin(ud.Coin),
		DelegatorAddress:    ud.DelegatorAddress,
		ValidatorOperAddr:   ud.ValidatorOperAddr,
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
	return model.StakingParams{
		Params: model.SParams{
			UnbondingTime:     sp.UnbondingTime,
			MaxValidators:     sp.MaxValidators,
			MaxEntries:        sp.MaxEntries,
			HistoricalEntries: sp.HistoricalEntries,
			BondDenom:         sp.BondDenom,
			MinCommissionRate: sp.MinCommissionRate.MustFloat64(),
		},
		Height: sp.Height,
	}
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
