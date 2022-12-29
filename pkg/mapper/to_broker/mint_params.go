package to_broker

import (
	"github.com/hexy-dev/spacebox/broker/model"

	"github.com/hexy-dev/spacebox-crawler/types"
)

func (tb ToBroker) MapMingParams(mp types.MintParams) model.MintParams {
	return model.NewMintParams(mp.Height, mp.MintDenom, mp.InflationRateChange.MustFloat64(),
		mp.InflationMax.MustFloat64(), mp.InflationMin.MustFloat64(), mp.GoalBonded.MustFloat64(), mp.BlocksPerYear)
}
