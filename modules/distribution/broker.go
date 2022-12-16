package distribution

import (
	"context"

	"github.com/hexy-dev/spacebox/broker/model"
)

type broker interface {
	PublishCommunityPool(ctx context.Context, cp model.CommunityPool) error
	PublishDistributionParams(ctx context.Context, dp model.DistributionParams) error
}
