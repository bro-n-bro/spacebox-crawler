package to_broker

import (
	"github.com/hexy-dev/spacebox-crawler/types"

	"github.com/hexy-dev/spacebox/broker/model"
)

func (tb ToBroker) MapValidatorInfo(val types.StakingValidator) model.ValidatorInfo {
	info := model.ValidatorInfo{
		ConsensusAddress:    val.GetConsAddr(),
		OperatorAddress:     val.GetOperator(),
		SelfDelegateAddress: val.GetSelfDelegateAddress(),
		Height:              val.GetHeight(),
	}

	if val.GetMinSelfDelegation() != nil {
		info.MinSelfDelegation = val.GetMinSelfDelegation().Int64()
	}

	return info
}
