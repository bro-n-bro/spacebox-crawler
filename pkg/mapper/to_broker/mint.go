package to_broker

import (
	"bro-n-bro-osmosis/adapter/broker/model"
	"bro-n-bro-osmosis/types"
)

func (tb ToBroker) MapMingParams(mp types.MintParams) model.MintParams {
	return model.MintParams{
		Height: mp.Height,
		Params: model.MParams{
			MintDenom:           mp.MintDenom,
			InflationRateChange: mp.InflationRateChange.MustFloat64(),
			InflationMax:        mp.InflationMax.MustFloat64(),
			InflationMin:        mp.InflationMin.MustFloat64(),
			GoalBonded:          mp.GoalBonded.MustFloat64(),
			BlocksPerYear:       mp.BlocksPerYear,
		},
	}
}
