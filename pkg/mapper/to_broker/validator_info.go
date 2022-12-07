package to_broker

import (
	"bro-n-bro-osmosis/adapter/broker/model"
	"bro-n-bro-osmosis/types"
)

func (tb ToBroker) MapValidatorInfo(val types.StakingValidator) model.ValidatorInfo {
	return model.ValidatorInfo{
		ConsensusAddress:    val.GetConsAddr(),
		OperatorAddress:     val.GetOperator(),
		SelfDelegateAddress: val.GetSelfDelegateAddress(),
		MinSelfDelegation:   val.GetMinSelfDelegation().Int64(),
		Height:              val.GetHeight(),
	}
}
