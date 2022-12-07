package to_broker

import (
	"bro-n-bro-osmosis/adapter/broker/model"
	"bro-n-bro-osmosis/types"
)

func (tb ToBroker) MapDistributionParams(params types.DistributionParams) model.DistributionParams {
	return model.DistributionParams{
		Height: params.Height,
		Params: model.DParams{
			CommunityTax:        params.Params.CommunityTax.MustFloat64(),
			BaseProposerReward:  params.Params.BaseProposerReward.MustFloat64(),
			BonusProposerReward: params.Params.BonusProposerReward.MustFloat64(),
			WithdrawAddrEnabled: params.Params.WithdrawAddrEnabled,
		},
	}
}
