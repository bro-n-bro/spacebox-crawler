package rank

import (
	"context"

	"github.com/bro-n-bro/spacebox/broker/model"
)

type broker interface {
	PublishRankParams(ctx context.Context, mp model.RankParams) error
}
