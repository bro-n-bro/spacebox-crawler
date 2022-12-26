package to_broker

import (
	"bro-n-bro-osmosis/types"

	"github.com/hexy-dev/spacebox/broker/model"
)

func (tb ToBroker) MapValidatorStatus(s types.ValidatorStatus) model.ValidatorStatus {
	return model.ValidatorStatus{
		Height:           s.Height,
		ValidatorAddress: s.ConsensusAddress,
		Status:           int64(s.Status),
		Jailed:           s.Jailed,
	}
}

func (tb ToBroker) MapValidatorsStatuses(statuses []types.ValidatorStatus) []model.ValidatorStatus {
	res := make([]model.ValidatorStatus, len(statuses))
	for i, status := range statuses {
		res[i] = tb.MapValidatorStatus(status)
	}
	return res
}
