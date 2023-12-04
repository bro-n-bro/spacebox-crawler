package dmn

import (
	"context"

	"github.com/bro-n-bro/spacebox/broker/model"
)

type broker interface {
	PublishDMNParams(ctx context.Context, mp model.DMNParams) error
}
