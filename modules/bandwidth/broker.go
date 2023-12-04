package bandwidth

import (
	"context"

	"github.com/bro-n-bro/spacebox/broker/model"
)

type broker interface {
	PublishBandwidthParams(ctx context.Context, mp model.BandwidthParams) error
}
