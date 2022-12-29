package to_broker

import (
	"github.com/hexy-dev/spacebox/broker/model"

	"github.com/hexy-dev/spacebox-crawler/types"
)

func (tb ToBroker) MapDistributionParams(params types.DistributionParams) model.DistributionParams {
	return model.NewDistributionParams(params.Height,
		params.CommunityTax.MustFloat64(),
		params.BaseProposerReward.MustFloat64(),
		params.BonusProposerReward.MustFloat64(),
		params.WithdrawAddrEnabled)
}
