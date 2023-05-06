package slashing

import (
	"context"

	"github.com/bro-n-bro/spacebox/broker/model"
)

type broker interface {
	PublishUnjailMessage(ctx context.Context, msg model.UnjailMessage) error
}
