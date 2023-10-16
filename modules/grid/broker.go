package grid

import (
	"context"

	"github.com/bro-n-bro/spacebox/broker/model"
)

type broker interface {
	PublishGridParams(ctx context.Context, mp model.GridParams) error
}
