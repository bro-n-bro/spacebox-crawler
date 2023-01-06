package mint

import (
	"context"

	"github.com/hexy-dev/spacebox/broker/model"
)

type broker interface {
	PublishMintParams(ctx context.Context, mp model.MintParams) error
	PublishAnnualProvision(ctx context.Context, ap model.AnnualProvision) error
}
