package to_broker

import (
	"bro-n-bro-osmosis/adapter/broker/model"
	"bro-n-bro-osmosis/types"
)

func (tb ToBroker) MapGovParams(params *types.GovParams) model.GovParams {
	return model.GovParams{
		DepositParams: model.DepositParams{
			MinDeposit:       tb.MapCoins(params.DepositParams.MinDeposit),
			MaxDepositPeriod: params.DepositParams.MaxDepositPeriod,
		},
		VotingParams: model.VotingParams{
			VotingPeriod: params.VotingParams.VotingPeriod,
		},
		TallyParams: model.TallyParams{
			Quorum:        params.TallyParams.Quorum.MustFloat64(),
			Threshold:     params.TallyParams.Threshold.MustFloat64(),
			VetoThreshold: params.TallyParams.VetoThreshold.MustFloat64(),
		},
		Height: params.Height,
	}
}
