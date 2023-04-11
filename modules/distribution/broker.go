package distribution

import (
	"context"

	"github.com/bro-n-bro/spacebox/broker/model"
)

type broker interface {
	PublishDelegationReward(context.Context, model.DelegationReward) error
	PublishDelegationRewardMessage(context.Context, model.DelegationRewardMessage) error
	PublishCommunityPool(ctx context.Context, cp model.CommunityPool) error
	PublishDistributionParams(ctx context.Context, dp model.DistributionParams) error
	PublishSetWithdrawAddressMessage(ctx context.Context, message model.SetWithdrawAddressMessage) error
	PublishProposerReward(ctx context.Context, pr model.ProposerReward) error
	PublishDistributionCommission(ctx context.Context, commission model.DistributionCommission) error
	PublishDistributionReward(ctx context.Context, reward model.DistributionReward) error
}
